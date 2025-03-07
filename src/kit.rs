struct Simulation {
    entities: Vec<Entity>,
    timeline: Vec<Event>,
    multiverse: Vec<Universe>,
}

struct Entity {
    id: char,
    states: Vec<String>, // Entities exist in multiple states across universes
}

struct Event {
    timestamp: u64,
    entity_id: char,
    action: String,
    universe_id: usize, // Tracks which universe this event belongs to
}

struct Universe {
    id: usize,
    observers: Vec<char>, // Entities that define what is real in this universe
}

fn main() {
    let entities = vec![
        Entity { id: '👻', states: vec!["".to_string()] },
        Entity { id: '🚁', states: vec!["🔼💨⏳".to_string()] },
        Entity { id: '🌞', states: vec!["🌚".to_string()] },
        Entity { id: '🦠', states: vec!["🦠🌝".to_string()] },
        Entity { id: '🏙️', states: vec!["🏙️".to_string()] },
        Entity { id: '🛰️', states: vec!["📡🔄".to_string()] },
    ];

    let mut simulation = Simulation {
        entities: entities.clone(),
        timeline: vec![],
        multiverse: vec![Universe {
            id: 0,
            observers: vec!['👻'],
        }],
    };
    
    run_simulation(&mut simulation, &entities);
    append_to_source();
    print_source();
}

fn run_simulation(sim: &mut Simulation, entities: &Vec<Entity>) {
    for universe in &mut sim.multiverse {
        for i in 0..entities.len() {
            if i > 0 && is_collision(&entities[i - 1], &entities[i], universe) {
                println!("⛔ Timeline disturbance detected in universe {}! Collision between {} and {}!", 
                         universe.id, entities[i - 1].id, entities[i].id);
                println!("🔍 Causal agent identified: {}", entities[i - 1].id);
                if universe.observers.contains(&'👻') {
                    branch_universe(sim, universe.id, entities[i - 1].id);
                }
            }
        }
        if universe.observers.contains(&'👻') {
            execute_sun_event(sim, universe.id);
        }
    }
}

fn is_collision(entity1: &Entity, entity2: &Entity, universe: &Universe) -> bool {
    let non_matter_entities = vec!['👻']; // Define non-material entities
    let entity2_is_matter = !non_matter_entities.contains(&entity2.id); // All else is matter
    let collision = entity1.id == '🚁' && entity2_is_matter; // Helicopter collides with matter
    
    if collision && universe.observers.contains(&'👻') {
        return true;
    }
    false
}

fn branch_universe(sim: &mut Simulation, parent_id: usize, cause: char) {
    let new_id = sim.multiverse.len();
    let new_universe = Universe {
        id: new_id,
        observers: vec![cause, '👻'], // The cause and original observer persist
    };
    println!("🌌 Branching new universe {} due to {}", new_id, cause);
    sim.multiverse.push(new_universe);
}

fn execute_sun_event(sim: &mut Simulation, universe_id: usize) {
    println!("🌞 Event triggered in universe {} by observer 👻", universe_id);
    for entity in &mut sim.entities {
        if entity.id == '🌞' {
            entity.states.push("🌚".to_string()); // Modify sun's state
            println!("🌚 The sun fades!");
        }
    }
}

fn append_to_source() {
    let filename = file!();
    let additional_line = "// 🐍 Self-replicating entity evolves\n";
    std::fs::OpenOptions::new()
        .append(true)
        .open(filename)
        .and_then(|mut file| std::io::Write::write_all(&mut file, additional_line.as_bytes()))
        .expect("Failed to append to source code");
}

fn print_source() {
    let source = std::fs::read_to_string(file!()).expect("Failed to read source code");
    println!("\nQuine Output:\n\n{}", source);
}
