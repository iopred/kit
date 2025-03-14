// use crate::kit::KitValue;

mod kv;

use kv::KitValue;

struct Simulation {
    entities: Vec<Entity>,
    timeline: Vec<Event>,
    multiverse: Vec<Universe>,
    observers: Vec<KitValue>, // Entities that can observe the simulation.
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

pub(crate) fn kit(input: &str) -> KitValue {
    let mut entities = vec![
        Entity { id: KitValue::String("👻".to_string()), states: vec![KitValue::String("🌌".to_string())] },
        Entity { id: KitValue::String("🔍".to_string()), states: vec![KitValue::String("🧑".to_string())] },
        Entity { id: KitValue::String("🚁".to_string()), states: vec![KitValue::String("🔼💨⏳".to_string())] },
        Entity { id: KitValue::String("🌞".to_string()), states: vec![KitValue::String("🌚".to_string())] },
        Entity { id: KitValue::String("🦠".to_string()), states: vec![KitValue::String("🦠🌝".to_string())] },
        Entity { id: KitValue::String("🏙️".to_string()), states: vec![KitValue::String("🏙️".to_string())] },
        Entity { id: KitValue::String("🛰️".to_string()), states: vec![KitValue::String("📡🔄".to_string())] },
    ];

    if !input.is_empty() {
        let id = input.chars().next().unwrap().to_string();
        let state = input.chars().skip(1).collect::<String>();
        entities.push(Entity {
            id: KitValue::String(id),
            states: vec![KitValue::String(state)],
        });
    }

    let mut simulation = Simulation {
        entities: entities.clone(),
        timeline: vec![],
        multiverse: vec![Universe {
            id: KitValue::Byte(0),
            observers: vec![KitValue::String("👻".to_string())],
        }],
        observers: vec![],
    };
    
    run_simulation(&mut simulation, &entities, 1);
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

fn branch_universe(sim: &mut Simulation, _parent_id: KitValue, cause: KitValue) {
    let new_id = sim.multiverse.len() as f64;
    let new_universe = Universe {
        id: KitValue::Boolean(true), // A new universe must always be contained in universe 1, otherwise our idea doesn't hold.
        observers: vec![cause.clone(), KitValue::String("👻".to_string())], // The cause and original observer persist
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

fn random_radius(radius: i32, center_x: f64, center_y: f64) -> (f64, f64) {
    let r = radius as f64 * (rand::random::<f64>().sqrt());
    let theta = rand::random::<f64>() * 2.0 * std::f64::consts::PI;
    let x = center_x + r * theta.cos();
    let y = center_y + r * theta.sin();
    (x, y)
}

fn append_to_source(sim: &mut Simulation) {
    let filename = file!();
    println!("boop{}", filename);
    let additional_line = format!("// {} Self-replicating entity evolves [🌞]\n", execute_sun_event(sim, KitValue::Boolean(false)));
    std::fs::OpenOptions::new()
        .append(true)
        .open(filename)
        .and_then(|mut file| std::io::Write::write_all(&mut file, additional_line.as_bytes()))
        .expect("Failed to append to source code");
}

fn print_source() -> String {
    let source = std::fs::read_to_string(file!()).expect("Failed to read source code");
    println!("{}", source);
    source
}
//  Self-replicating entity evolves [🌞]
//  Self-replicating entity evolves [🌞]
//  Self-replicating entity evolves [🌞]
//  Self-replicating entity evolves [🌞]
