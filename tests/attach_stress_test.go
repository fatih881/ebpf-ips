//go:build stress

package tests

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	localebpf "github.com/fatih881/ebpf-ips/core/ebpf"
	"github.com/vishvananda/netlink"
)

const testBinary = "stressTest"
const ifacesBeforeRun = 50
const ifacesAfterRun = 50
const ifacesToDelete = 30
const ifacesToLinkFlap = 20
const durationLinkFlap = 10 // in seconds
const waitTime = 10         // for prometheus to return correct info,in seconds
const waitHTTP = 10         // for main.go to start listening 2112(prometheus),in seconds
func TestAttachStress(t *testing.T) {
	attachStress(t)
}
func attachStress(t *testing.T) {
	objs := &localebpf.IpsObjects{}
	err := localebpf.LoadIpsObjects(objs, nil)
	if err != nil {
		t.Fatalf("cannot load program to kernel : %v", err)
	}
	t.Cleanup(func() {
		if err := objs.Close(); err != nil {
			t.Logf("failed to remove objects : %v", err)
		}
	})
	build := exec.Command("go", "build", "-a", "-v", "-o", testBinary, "../main.go")
	build.Env = os.Environ()
	build.Stdout = os.Stdout
	build.Stderr = os.Stderr
	if err := build.Run(); err != nil {
		t.Fatalf("failed to build program : %v", err)
	}
	var dummyLinks []netlink.Link
	var kernelIndexes []int
	var availableIfaceNumber int
	availableIfaces, err := net.Interfaces()
	if err != nil {
		t.Fatalf("cannot get available interfaces : %v", err)
	}
	for _, iface := range availableIfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		availableIfaceNumber++
	}
	if err := bulkCreateIfaces(t, 1, ifacesBeforeRun, &dummyLinks, &kernelIndexes); err != nil {
		t.Fatalf("error creating ifaces : %v", err)
	}
	chmod := exec.Command("chmod", "+x", testBinary)
	chmod.Stdout = os.Stdout
	chmod.Stderr = os.Stderr
	if err := chmod.Run(); err != nil {
		t.Fatalf("error grant permission : %v", err)
	}
	runbinary := exec.Command("./" + testBinary)
	runbinary.Stdout = os.Stdout
	runbinary.Stderr = os.Stderr
	if err := runbinary.Start(); err != nil {
		t.Fatalf("error running binary : %v", err)
	}
	t.Cleanup(func() {
		if err := runbinary.Process.Kill(); err != nil {
			t.Logf("failed to kill binary : %v", err)
		}
	})
	if err := bulkCreateIfaces(t, ifacesBeforeRun+1, ifacesBeforeRun+ifacesAfterRun, &dummyLinks, &kernelIndexes); err != nil {
		t.Fatalf("error creating ifaces : %v", err)
	}
	if err := bulkReadXdpProgs(&kernelIndexes); err != nil {
		t.Logf("error(kernel) XDP program could not be attached to all test interfaces: %v", err)
	}
	if err := waitForPrometheusServer("http://localhost:2112/metrics", waitHTTP*time.Second); err != nil {
		t.Fatalf("cannot reach prometheus after %v : %v", waitHTTP, err)
	}
	resp, err := http.Get("http://localhost:2112/metrics")
	if err != nil {
		t.Fatalf("cannot get prometheus metrics: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("prometheus returned : %d", resp.StatusCode)
	}
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close body : %v", err)
		}
	})
	metrics, err := io.ReadAll(resp.Body)
	t.Logf("Waiting for attach success: %d", availableIfaceNumber+ifacesBeforeRun+ifacesAfterRun)
	if err := waitForMetric(t, "http://localhost:2112/metrics", "xdp_attach_success_total", availableIfaceNumber+ifacesBeforeRun+ifacesAfterRun, waitTime*time.Second); err != nil {
		t.Fatalf("Attach metrics check failed: %v", err)
	}
	if err := bulkDeleteIfaces(1, ifacesToDelete, &dummyLinks); err != nil {
		t.Logf("error removing ifaces : %v", err)
	}
	kernelIndexes = kernelIndexes[ifacesToDelete:]
	t.Logf("Waiting for detach success: %d", ifacesToDelete)
	if err := waitForMetric(t, "http://localhost:2112/metrics", "xdp_detach_total", ifacesToDelete, waitTime*time.Second); err != nil {
		t.Fatalf("Detach metrics check failed: %v", err)
	}
	resp, err = http.Get("http://localhost:2112/metrics")
	if err != nil {
		t.Fatalf("cannot get prometheus metrics: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("prometheus returned : %d", resp.StatusCode)
	}
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close body : %v", err)
		}
	})
	metrics, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("cannot read prometheus metrics: %v", err)
	}
	if err := comparePrometheusResult(ifacesBeforeRun+ifacesAfterRun+availableIfaceNumber, 0, ifacesToDelete, metrics); err != nil {
		t.Logf("error comparing prometheus result : %v", err)
	}
	if err := bulkReadXdpProgs(&kernelIndexes); err != nil {
		t.Logf("error(kernel) XDP program could not be attached to all test interfaces: %v", err)
	}
	bulkLinkFlap(t, ifacesToDelete+1, ifacesToDelete+ifacesToLinkFlap, &dummyLinks, durationLinkFlap*time.Second)
	if err := comparePrometheusResult(ifacesAfterRun+availableIfaceNumber, 0, ifacesToDelete, metrics); err != nil {
		t.Logf("error comparing prometheus result : %v", err)
	}
	if err := bulkReadXdpProgs(&kernelIndexes); err != nil {
		t.Logf("error(kernel) XDP program could not be attached to all test interfaces: %v", err)
	}
	resp, err = http.Get("http://localhost:2112/metrics")
	if err != nil {
		t.Fatalf("cannot get prometheus metrics: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("prometheus returned : %d", resp.StatusCode)
	}
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("failed to close body : %v", err)
		}
	})
	metrics, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("cannot read prometheus metrics: %v", err)
	}
	if err := comparePrometheusResult(ifacesBeforeRun+ifacesAfterRun+availableIfaceNumber, 0, ifacesToDelete, metrics); err != nil {
		t.Logf("error comparing prometheus result : %v", err)
	}
}
func bulkCreateIfaces(t *testing.T, startIfaceNum int, endIfaceNum int, dummyLinks *[]netlink.Link, kernelIndexes *[]int) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := startIfaceNum; i <= endIfaceNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			dummyLinkName := fmt.Sprintf("xdp-e2e-%03d", id)
			dummyLink := &netlink.Dummy{
				LinkAttrs: netlink.LinkAttrs{Name: dummyLinkName},
			}
			if err := netlink.LinkAdd(dummyLink); err != nil {
				t.Logf("failed to add link %s: %v", dummyLinkName, err)
				return
			}

			if err := netlink.LinkSetUp(dummyLink); err != nil {
				return
			}
			info, _ := net.InterfaceByName(dummyLinkName)
			mu.Lock()
			*kernelIndexes = append(*kernelIndexes, info.Index)
			*dummyLinks = append(*dummyLinks, dummyLink)
			mu.Unlock()

		}(i)
	}

	wg.Wait()
	return nil
}
func bulkDeleteIfaces(startIfaceNum int, endIfaceNum int, dummyLinks *[]netlink.Link) error {
	for i := startIfaceNum; i <= endIfaceNum; i++ {
		dummyLinkName := fmt.Sprintf("xdp-e2e-%03d", i)
		var linkToDelete netlink.Link
		for _, link := range *dummyLinks {
			if link.Attrs().Name == dummyLinkName {
				linkToDelete = link
				break
			}
		}
		if linkToDelete == nil {
			return fmt.Errorf("link %s not found in dummyLinks", dummyLinkName)
		}
		err := netlink.LinkDel(linkToDelete)
		if err != nil {
			return fmt.Errorf("cannot delete dummy for xdp : %v", err)
		}
	}
	return nil
}

func bulkReadXdpProgs(kernelIndexes *[]int) error {
	var missingXDP []int
	for i := 0; i < len(*kernelIndexes); i++ {
		linkObj, err := netlink.LinkByIndex((*kernelIndexes)[i])
		if err != nil {
			return fmt.Errorf("error at query : %v", err)
		}
		if linkObj.Attrs().Xdp == nil {
			missingXDP = append(missingXDP, (*kernelIndexes)[i])
		}
	}
	if len(missingXDP) > 0 {
		return fmt.Errorf("XDP not attached to %d interfaces: %v", len(missingXDP), missingXDP)
	}
	return nil
}
func comparePrometheusResult(success int, fail int, detach int, metrics []byte) error {
	attachSuccess := fmt.Sprintf("xdp_attach_success_total %d", success)
	if !strings.Contains(string(metrics), attachSuccess) {
		return fmt.Errorf("expected %s , full body %s", attachSuccess, string(metrics))
	}
	attachFail := fmt.Sprintf("xdp_attach_failed_total %d", fail)
	if !strings.Contains(string(metrics), attachFail) {
		return fmt.Errorf("expected %s , full body %s", attachFail, string(metrics))
	}
	detachTotal := fmt.Sprintf("xdp_detach_total %d", detach)
	if !strings.Contains(string(metrics), detachTotal) {
		return fmt.Errorf("expected %s , full body %s", detachTotal, string(metrics))
	}
	return nil
}
func waitForMetric(t *testing.T, url string, fetchedMetric string, expectedMetric int, timeout time.Duration) error {
	var body []byte
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for time.Now().Before(deadline) {
		select {
		case <-ticker.C:
			resp, err := http.Get(url)
			if err != nil {
				continue
			}
			body, _ = io.ReadAll(resp.Body)
			err = resp.Body.Close()
			if err != nil {
				t.Logf("error closing response body: %v", err)
				return err
			}
			if strings.Contains(string(body), fmt.Sprintf("%s %d", fetchedMetric, expectedMetric)) {
				return nil
			}
		}
	}
	return fmt.Errorf("timeout: metric %s did not reach %d within %v. Last body: %s", fetchedMetric, expectedMetric, timeout, string(body))
}
func bulkLinkFlap(t *testing.T, startNum int, endNum int, dummyLinks *[]netlink.Link, duration time.Duration) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	for i := startNum; i <= endNum; i++ {
		wg.Add(1)
		go func(currentIndex int) {
			defer wg.Done()
			dummyLinkName := fmt.Sprintf("xdp-e2e-%03d", currentIndex)
			var linkToFlap netlink.Link
			for _, link := range *dummyLinks {
				if link.Attrs().Name == dummyLinkName {
					linkToFlap = link
					break
				}
			}
			if linkToFlap == nil {
				t.Logf("could not find link %s(flap)", dummyLinkName)
				return
			}
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if rand.Intn(2) == 0 {
						if err := netlink.LinkSetDown(linkToFlap); err != nil {
							t.Logf("error at setting link down(flap) : %v", err)
						}
					} else {
						if err := netlink.LinkSetUp(linkToFlap); err != nil {
							t.Logf("error at setting link down(flap) : %v", err)
						}
					}
				}
				if err := netlink.LinkSetUp(linkToFlap); err != nil {
					t.Logf("error at setting link up after stress test(flap) : %v", err)
				}
				time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
			}
		}(i)
	}
	wg.Wait()
}
func waitForPrometheusServer(url string, duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("no reply from %s : timed out in %v", url, duration)
		default:
			resp, err := http.Get(url)
			if err == nil {
				err := resp.Body.Close()
				if err != nil {
					return err
				}
				if resp.StatusCode == 200 {
					return nil
				}
			}

		}
		time.Sleep(500 * time.Millisecond)
	}
}
