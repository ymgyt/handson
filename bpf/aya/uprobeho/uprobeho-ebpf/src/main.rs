#![no_std]
#![no_main]

use aya_ebpf::{macros::uprobe, programs::ProbeContext};
use aya_log_ebpf::info;

#[uprobe]
pub fn uprobeho(ctx: ProbeContext) -> u32 {
    match try_uprobeho(ctx) {
        Ok(ret) => ret,
        Err(ret) => ret,
    }
}

fn try_uprobeho(ctx: ProbeContext) -> Result<u32, u32> {
    info!(&ctx, "function uprobed_function called by /proc/self/exe");
    let arg: u32 = ctx.arg(0).unwrap_or(9999);
    info!(&ctx, "arg: {}", arg);
    Ok(0)
}

#[cfg(not(test))]
#[panic_handler]
fn panic(_info: &core::panic::PanicInfo) -> ! {
    loop {}
}
