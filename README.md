### Welcome!
Developing game: mechanics are based on my novels. Trial number four! 

# EXPERIMENT FAILED
- Can't publish stats behind latest state of prev stats,
- So can't publish fetch between states (only latest state).

Need total rework.  
Meet the new trial MAG-Epsilon soon.  

# Usage
 - For build with go v.18+  
   `go build .`  
 - For run with go v.18+  
 `go run . -p 1234-123456789 -d cache/bc`  

Command:
```bash
mag-delta 

    -h # for help
    # Default: false

    -p 0000-000000000 # player ID to login (insecure yet), 
    # to be deprecated with built-in CLI client
    # Default: [not defined]

    -n # to generate new player, 
    #to be deprecated with built-in CLI client
    # Default: false

    -d ./path # blockchain work directory, automatically created 
    # Default: ./cache

```

# Gameplay
Dummy is getting spawned in front of you.
Destroy it as soon as possible!

Skills:  
- `+` to update your character stats  
  *Useless for now, because it improves nothing except last time*
- `e` to jinx target  
  *Jinx - simle attacking skill with your energy, element: none (from the stream)*

Menu: 
- `/` to list chain meta
- `?` to list parent blocks

# Versions
- `v0...`: Dummies
  - `v0.0...` Player chain
    - `v0.0.1`: NIF-N33 Dummies slayer: Wrestless wood  
      Jinx the wood! 
    - `v0.0.2`: N33-N3G **Dummies slayer: natural phenomena**  
      **Jinx the targed with all you might!**
      - Compressed objects (floats to int, structs to maps).  
      - Multiple streams per playable character.  
  - `v0.1...` *P2P* or *Spawn chain*
  - `v0.2...` *Spawn chain* or *P2P*

# Data scheme 
`/` - Root, genesis block
- `/Players` - last player created, initial player blocks (born blocks)
  - `/Players/0000-000000000` - stats for user id 0000-000000000
    - `/Session/0000-000000000/0-0000000` - session for player 0000-000000000 with stats 0-0000000  
    - `/Session/0000-000000000` - latest session for player

# Primitives basics:
|Object|Properties|Description|
|---:|:---:|:---|
|Player|Pool, Stream, Health|That's you, me dear friend|
|Health|Max, Current|obviously, the most valuable resource|
|Stream|Length (Creation), Width (Alteration), Power (Destruction), Nature (Element)|Universal vector to define any natural phenomena|
|Pool|Max, Dots|Container for unused dots, |
|Dot|Weight, Element|Energy point, most usable resource|


### Phenomena basics:
Everything is natural phenomena!

|v Minor / Major >|Cre|Alt|Des|
|---:|:---:|:---:|:---:|
|**Cre**ation|Conjuration|Boon|Barrier|
|**Alt**eration|Restoration|Modification|Enchantment|
|**Des**truction|Decay|Condition|Damage|
|All|Summoning|Transfiguration|Disaster|

### Elements:

|Rarity|Element|Cre|Alt|Des|
|---:|:---:|---|---|---|
||**✳️Common**||||
|1|☁️Air|Pressure* (Accumulation)|Spreading (AOE, Distance)|Penetration|
|1|🔥Fire|Energy* (Fueling)|Warming (Heat)|Burn|
|1+|⛰Earth|Structure|Mass*|Grind|
|2|❄️Water|Form|Direction|Dissolution*|
||**✳️✳️ Superior**||||
|3|✳️✳️[+✳️] Composed commons|Major's cre|Major's alt|Major's des|
|4|🌑Void|Shadows|Curse|Pain|
|4+|🌑✳️ Composed with void|Minor's cre|Minor's alt|Minor's des|
||**✳️✳️✳️ Deviant**||||
|5|🩸Mallom|Grow|Corruption|Consumption|
|5+|🎶Noise|Inspiration|Echo (Confusion)|Roar (Fear)|
|7|🌟Light|Mirage|Matter|Desintegration|
