//go:build ignore 
#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/in.h>
#include <bpf/bpf_helpers.h>
// send key,get value as return.
struct {
    __uint(type, BPF_MAP_TYPE_HASH); 
    __uint(max_entries, 100);           
    __type(key, __u32);               
    __type(value, __u32);             
} allowed_ip_map SEC(".maps");   
struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(max_entries, 100);
    __type(key, __u8);
    __type(value, __u8);
} allowed_ports SEC(".maps"); 
// ALLOWING PACKETS VIA PORTS IS CURRENTLY UNAVAILABLE
SEC("xdp")
int ips(struct xdp_md *ctx){
    void *data =(void *)(long)ctx->data; // data start.
    void *data_end =(void *)(long)ctx->data_end;
    struct ethhdr *eth = data;
    struct iphdr *ip;
    // This if checks for ** is header is in the packet data ** or not.
    // Example : If the pointer for data_start (named as data in xdp_md) is in the 100. ram block,and the header is 20 bytes,the ethernet header will be
    // ended at 120.If the packet is ended at 119,this packet is corrupted.    
    if ((void *)(eth + 1) > data_end)
    // I dropped the packet,because it will be dropped even if I'll pass it.
    // You can XDP_PASS if you want to see it in higher-level monitoring,as like tcpdump.
    // This programs can't see XDP_DROP's because it's dropped before it reaches kernel.
        return XDP_DROP;
    //This if checks for the protocol of eth(in struct,it's raw data)is ipv4 or not.
    // htonl stands for host to network long,because x86 reverses the raw data,we reverse it again with htonl.
    if (eth->h_proto != __constant_htons(ETH_P_IP)){
        return XDP_PASS;
        // ip_key = ip->saddr;
        // ports_key = ip ->protocol;
    } 
    //After this check,we will don't touch to other protocols such as ARP or ipv6.I will may add ipv6's.
    //
    ip = (struct iphdr *)(eth + 1);  // as we know it's a ipv4 package,we can define it now.
    if ((void *)(ip + 1) > data_end)
        return XDP_DROP;
    //
    __u32 *allowed_source = bpf_map_lookup_elem(&allowed_ip_map, &ip->saddr);
   /* __u8  *allowed_port   = bpf_map_lookup_elem(&allowed_ports, &ip->protocol);
    if (!allowed_source || !allowed_port)
        return XDP_DROP;
        */
    if (!allowed_source)
        return XDP_DROP;
    return XDP_PASS;
}
char LICENSE[] SEC("license") = "GPL";