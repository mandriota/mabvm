# Mab Virtual Machine
**Stack Virtual Machine named after Mab - Queen of The Faires**
```text
                        . _______      . _________
     . _____      . _____.|     |      |.|       |__________________________
    / /     \    / /     \.\    \_______| |     |.|      _____________ _     \
   / /   /\  \  / /   /\  \.\     ________      |.|      |____________\ \     |
  / /   /. \  \/ /   /. \  \.\    \     | |     |.|      ________________    <
 / /   /  . \  \/   /  . \  \.\    \____|.|     |.|      |____________/_/     |
|.|____|  |.|_______| |.|______|________________|.|__________________________/

            .  .  .
            |\/ \/|
           .~+/~_``\.
         ,// '   '\||\`
        '/||<*)  (`\\\\|`
        /|||   /  `|\\||`
       ,/// \  ~   /`\\||`
       '|||  |`  `|  ||||'
       '|||_'|    |'_||||'
       /'   '       '   '\
```

## Memory Model
A linear array divided in 8x4KiB blocks.
Every single block have only one owner indicated with mutex.
If a proccess wants to write not own block, block mutex must be locked.

## Opcode Model
First 2 bits describes sequence.
Other 6 bits describes control and conditional flags.
| JUMP  | CONTROL FLAGS | CONDITIONAL FLAGS |
| ----- | ------------- | ----------------- |
| 2-bit | 3-bit         | 3-bit             |

#### Assembler Instruction
Each instruction is divided in 3 sections:
- Jump-Sequence Section, starts with `:`
- Control Flags Section, starts with `'`
- Conditional Flags Section, starts with `"`

##### Syntax
```
:<S|D|C|V>'[I][E][M]"[L][E][G]
```

##### Sequence-Jump Section
| code | letter | src (`=dst−k`) | dst (`=src+k`) |
| ---- | ------ | -------------- | -------------- |
| 0x0  | S      | data src ptr   | data src ptr   |
| 0x1  | D      | data dst ptr   | data dst ptr   |
| 0x2  | C      | code ptr       | code ptr       |
| 0x3  | V      | data src       | data dst       |

##### Control-Flags Section
| code | letter | action                      | description                |
| ---- | ------ | --------------------------- | -------------------------- |
| 0x4  | I      | `k=−k`                      | inverts `k`                |
| 0x8  | E      | `k=k×data[src]; src=src−1`  | extends `k`                |
| 0x10 | M      | `sleep for locked mutex`    | locks proccess until write |
*Indicates mode to execute instruction with 1-3 ordered characters.*

##### Conditional-Flags Section
| code | letter | meaning | action      | description                        |
| ---- | ------ | ------- | ----------- | ---------------------------------- |
| 0x20 | L      | lower   | `src<dst`   | source is lower than destination   |
| 0x40 | E      | equal   | `src=dst`   | source is equal to destination     |
| 0x80 | G      | greater | `src>dst`   | source is greater than destination |
*Indicates condition to execute instruction with 1-3 ordered characters.*
