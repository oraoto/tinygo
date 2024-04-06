MEMORY
{
	FLASH (rx) : ORIGIN = 0x00000000, LENGTH = 16K
	RAM (xrw)  : ORIGIN = 0x20000000, LENGTH = 2K
}

ENTRY(_start)

SECTIONS
{
    .init :
    {
        . = ALIGN(4);
        KEEP(*(SORT_NONE(.init)))
        . = ALIGN(4);
    } >FLASH AT>FLASH

    .vector_table :
    {
        . = ALIGN(4);
        KEEP(*(.vector_table.interrupts));
        . = ALIGN(4);
    } >FLASH AT>FLASH

    .text :
    {
        . = ALIGN(4);
        KEEP(*(SORT_NONE(.handle_reset)))
        *(.text .text.*)
    } >FLASH AT>FLASH

    .rodata : ALIGN(4)
    {
        *(.srodata .srodata.*);
        *(.rodata .rodata.*);
    } >FLASH AT>FLASH

    .stack (NOLOAD) :
    {
        . = ALIGN(4);
        . += _stack_size;
        _stack_top = .;
    } > RAM

    .data : ALIGN(4)
    {
        _data_lma = LOADADDR(.data);
        PROVIDE(_data_vma = .);
        _sdata = .;  /* used by startup code */
        *(.data .data.*)
        . = ALIGN(4);
        PROVIDE( __global_pointer$ = . + 0x800 );
        /* These sections are used by the BLE lib */
        *(.sdata .sdata.*)
        . = ALIGN(4);
        PROVIDE( _edata = .); /& used by startup code */
    } >RAM AT>FLASH

    .bss : ALIGN(4)
    {
        PROVIDE( _sbss = .);
        *(.sbss .sbss.*)
        *(.bss .bss.*)
        *(.gnu.linkonce.sb.*)
        *(.gnu.linkonce.b.*)
        *(COMMON*)
        PROVIDE( _ebss = .);
    } >RAM AT>FLASH

/*
    .stack ORIGIN(RAM)+LENGTH(RAM) :
    {
        . = ALIGN(4);
        PROVIDE(_stack_top = . );
    } >RAM


    .got (INFO) :
    {
        KEEP(*(.got .got.*));
    }

    .eh_frame (INFO) : { KEEP(*(.eh_frame)) }
    .eh_frame_hdr (INFO) : { *(.eh_frame_hdr) }
*/
}

/* For the memory allocator. */
_heap_start = ALIGN(_ebss, 4);
_heap_end = ORIGIN(RAM) + LENGTH(RAM);
_globals_start = _sdata;
_globals_end = _ebss;

_stack_size = 256;
