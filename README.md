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
Indicates with 1 character sequence to jump:
- `S` indicates Source-Value Jump-Sequence
- `D` indicates Destination-Value Jump-Sequence
- `C` indicated Counter-Value Jump-Sequence
- `V` indicated Value-Value Jump-Sequence

##### Control-Flags Section
Indicates control flags with 1-3 ordered characters:
1. `I` inverts value
2. `E` extends value
3. `M` locks proccess until one of blocks will rewrited

##### Conditional-Flags Section
Indicates instruction condition to executed with 1-3 ordered characters:
1. `L` executes if source is lower than destination
2. `E` executes if source is equal to destination
3. `G` executes if source is greater than destination