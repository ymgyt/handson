use rbpf::assembler::assemble;
use rbpf::helpers;

fn main() {}

#[test]
fn test_vm_add() {
    let prog = assemble(
        "
            mov32 r0, 1
            mov32 r1, 3
            add r0, r1
            exit",
    )
    .unwrap();

    let vm = rbpf::EbpfVmNoData::new(Some(&prog)).unwrap();
    let ret = vm.execute_program().unwrap();
    assert_eq!(ret, 0x4);
}
