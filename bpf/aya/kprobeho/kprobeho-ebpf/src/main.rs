#![no_std]
#![no_main]

#[allow(non_upper_case_globals)]
#[allow(non_snake_case)]
#[allow(non_camel_case_types)]
#[allow(dead_code)]
mod vmlinux;

use crate::vmlinux::{sock, sock_common};

use aya_ebpf::{macros::kprobe, programs::ProbeContext};
use aya_log_ebpf::info;

const AF_INET: u16 = 2;
const AF_INET6: u16 = 10;

#[kprobe]
pub fn kprobeho(ctx: ProbeContext) -> u32 {
    match try_kprobeho(ctx) {
        Ok(ret) => ret,
        Err(ret) => match ret.try_into() {
            Ok(rt) => rt,
            Err(_) => 1,
        },
    }
}

fn try_kprobeho(ctx: ProbeContext) -> Result<u32, i64> {
    let _sock: *mut sock = ctx.arg(0).ok_or(1i64)?;
    Ok(1)
}

#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    unsafe { core::hint::unreachable_unchecked() }
}
