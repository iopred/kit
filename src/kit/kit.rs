mod kv;

use kv::KitValue;

struct Simulation {
    entities: Vec<Entity>,
    timeline: Vec<Event>,
    multiverse: Vec<Universe>,
}

#[derive(Clone, Debug)]
struct Entity {
    id: KitValue,
    states: Vec<KitValue>, // Entities exist in multiple states across universes
}

struct Event {
    timestamp: KitValue,
    entity_id: KitValue,
    action: KitValue,
    universe_id: KitValue, // Tracks which universe this event belongs to
}

struct Universe {
    id: KitValue,
    observers: Vec<KitValue>, // Entities that can observe this universe
}

pub(crate) fn kit() -> KitValue {
    let entities = vec![
        Entity { id: KitValue::String("👻".to_string()), states: vec![KitValue::String("🌌".to_string())] },
        Entity { id: KitValue::String("🚁".to_string()), states: vec![KitValue::String("🔼💨⏳".to_string())] },
        Entity { id: KitValue::String("🌞".to_string()), states: vec![KitValue::String("🌚".to_string())] },
        Entity { id: KitValue::String("🦠".to_string()), states: vec![KitValue::String("🦠🌝".to_string())] },
        Entity { id: KitValue::String("🏙️".to_string()), states: vec![KitValue::String("🏙️".to_string())] },
        Entity { id: KitValue::String("🛰️".to_string()), states: vec![KitValue::String("📡🔄".to_string())] },
    ];

    let mut simulation = Simulation {
        entities: entities.clone(),
        timeline: vec![],
        multiverse: vec![Universe {
            id: KitValue::Number(0.0),
            observers: vec![KitValue::String("👻".to_string())],
        }],
    };
    
    run_simulation(&mut simulation, &entities, 60);
    append_to_source(&mut simulation);
    print_source();

    let mut kit_string = String::new();
    for entity in entities {
        kit_string.push_str(&format!("{:?}\n", entity));
    }

    KitValue::String(kit_string)
}

fn run_simulation(sim: &mut Simulation, entities: &Vec<Entity>, max_events: usize) {
    let mut event_count = 0;
    while event_count < max_events {
        let mut branches = vec![];
        let mut sun_events = vec![];

        for universe in &sim.multiverse {
            for i in 0..entities.len() {
                if i > 0 && is_collision(&entities[i - 1], &entities[i], universe) {
                    println!("⛔ Timeline disturbance detected in universe {:?}! Collision between {:?} and {:?}!", 
                             universe.id, entities[i - 1].id, entities[i].id);
                    println!("🔍 Causal agent identified: {:?}", entities[i - 1].id);
                    if universe.observers.contains(&KitValue::String("👻".to_string())) {
                        branches.push((universe.id.clone(), entities[i - 1].id.clone()));
                    }
                }
            }
            if universe.observers.contains(&KitValue::String("👻".to_string())) {
                sun_events.push(universe.id.clone());
            }
        }

        for (universe_id, cause) in branches {
            branch_universe(sim, universe_id, cause);
        }

        for universe_id in sun_events {
            execute_sun_event(sim, universe_id);
        }

        event_count += 1;
    }
}

fn is_collision(entity1: &Entity, entity2: &Entity, universe: &Universe) -> bool {
    let non_matter_entities = vec![KitValue::String("👻".to_string())]; // Define non-material entities
    let entity2_is_matter = !non_matter_entities.contains(&entity2.id); // All else is matter
    let collision = entity1.id == KitValue::String("🚁".to_string()) && entity2_is_matter; // Helicopter collides with matter
    
    if collision && universe.observers.contains(&KitValue::String("👻".to_string())) {
        return true;
    }
    false
}

fn branch_universe(sim: &mut Simulation, parent_id: KitValue, cause: KitValue) {
    let new_id = sim.multiverse.len() as f64;
    let new_universe = Universe {
        id: KitValue::Number(new_id),
        observers: vec![cause, KitValue::String("👻".to_string())], // The cause and original observer persist
    };
    println!("🌌 Branching new universe {} due to {:?}", new_id, cause);
    sim.multiverse.push(new_universe);
}

fn execute_sun_event(sim: &mut Simulation, universe_id: KitValue) -> String {
    println!("🌞 Event triggered in universe {:?} by observer 👻", universe_id);
    for entity in &mut sim.entities {
        if entity.id == KitValue::String("🌞".to_string()) {
            entity.states.push(KitValue::String("🌚".to_string())); // Modify sun's state
            println!("🌚 The sun fades!");
        }
    }
    match sim.entities.last().unwrap().states.last().unwrap() {
        KitValue::String(s) => s.clone(),
        _ => "".to_string(),
    }
}

fn append_to_source(sim: &mut Simulation) {
    let filename = file!();
    let additional_line = format!("// {} Self-replicating entity evolves [🌞]", execute_sun_event(sim, KitValue::Number(0.0)));
    std::fs::OpenOptions::new()
        .append(true)
        .open(filename)
        .and_then(|mut file| std::io::Write::write_all(&mut file, additional_line.as_bytes()))
        .expect("Failed to append to source code");
}

fn print_source() -> String {
    let source = std::fs::read_to_string(file!()).expect("Failed to read source code");
    println!("\nQuine Output:\n\n{}", source);
    source
}