
let jsonData = require("./data.json")


for (index in jsonData) {
    let creature = jsonData[index]
    let newCreature = {...creature}
    if (creature.spells?.length > 0) {
        newCreature.spellNotes = creature.spells[0]
        newCreature.spells = {}
        newCreature.precombatSpells = []
        for (let i = 0; i < 12; i++) {
            
            if (typeof(creature.spells[i+1]) !== 'string' && creature.spells[i+1] !== undefined) {
                for (key in creature.spells[i+1]) {
                    newCreature.spells[i] = creature.spells[i+1][key]
                }
                
                if (newCreature.spells[i].includes("*")) {
                    let list = newCreature.spells[i].split(",")
                    for (j in list) {
                        let spell = list[j]
                        if (spell.includes("*")) {
                            newCreature.precombatSpells = [...newCreature.precombatSpells, spell.trim()]
                        }
                    }

                
                }

            } else if (i > 0 && creature.spells[i+1] !== undefined) {
                newCreature.beforeCombat = creature.spells[i+1]
            }
        }
        console.log(newCreature.name)
        console.log(newCreature.spellNotes)
        console.log(newCreature.spells)
        console.log(newCreature.precombatSpells)
    }

    jsonData[index] = newCreature
}


var fs = require('fs')

fs.writeFileSync("new_data.json", JSON.stringify(jsonData, null).replace(RegExp(/\*/gm),""), function(err) {
    if (err) {
        console.log(err);
    }
})


