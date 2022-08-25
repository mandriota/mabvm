# Mab Virtual Machine
**Esoteric Stack Language Reduced Instruction Set Virtual Machine**

```txt:logo.txt [2-8]
```

## Memory Model

## Opcode Model
| JUMP  | FLAGS |
| ----- | ----- |
| 2-bit | 6-bit |

#### Assembler Instruction
:<S|D|C|V>:\[I][E][M]:\[L][E][G]

## Jump List

| no | asm | target         |
| -- | --- | -------------- |
| 0  |  S  | source         |
| 1  |  D  | distrubution   |
| 2  |  C  | code           |
| 3  |  V  | value          |

## Flag List

| no | asm | synopsis |
| -- | --- | -------- |
| 0  |  I  | inverse  |
| 1  |  E  | extend   |
| 2  |  M  | mutex    |
| 3  |  L  | lower    |
| 4  |  E  | equal    |
| 5  |  G  | greater  |