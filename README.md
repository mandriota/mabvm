# Mab Virtual Machine
**Esoteric Stack Language Virtual Machine**
```text
                        . _______      . _________
     . _____      . _____.|     |      |.|       |__________________________
    / /     \    / /     \.\    \_______| |     |.|      _____________ _     \
   / /   /\  \  / /   /\  \.\     ________      |.|      |____________\ \     |
  / /   /. \  \/ /   /. \  \.\    \     | |     |.|      ________________    <
 / /   /  . \  \/   /  . \  \.\    \____|.|     |.|      |____________/_/     |
|.|____|  |.|_______| |.|______|________________|.|__________________________/
```

## Memory Model
Linear array divided in 4KB blocks.
Every single block have only one owner indicated with mutex.
If a proccess wants to write not own block, block mutex must be locked.

## Opcode Model
First 2 bits describes jump-table.
Other 6 bits describes control and conditional flags.
| JUMP  | FLAGS |
| ----- | ----- |
| 2-bit | 6-bit |

#### Assembler Instruction
Each instruction is divided in 3 sections:
- Jump-Table Section
- Control Flags Section
- Conditional Flags Section
Every section must start with `:`.

##### Syntax
```
:<S|D|C|V>:\[I][E][M]:\[L][E][G]
```

##### Jump-Table Section
Indicates with 1 character table to jump:
- `S` indicates Source-Value Jump-Table
- `D` indicates Destination-Value Jump-Table
- `C` indicated Counter-Value Jump-Table
- `V` indicated Value-Value Jump-Table

##### Control-Flags Section
Indicates control flags with 0-3 ordered characters:
1. `I` inverts value
2. `E` extends value
3. `M` locks proccess until one of blocks will rewrited

##### Conditional-Flags Section
Indicates instruction condition to executed with 0-3 ordered characters:
1. `L` executes if source is lower than destination
2. `E` executes if source is equal to destination
3. `G` executes if source is greater than destination