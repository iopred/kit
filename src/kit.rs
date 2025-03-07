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
    entities: Vec<Entity>,
    observers: Vec<char>, // Entities that define what is real in this universe
}

fn main() {
    let mut simulation = Simulation {
        entities: vec![
            Entity { id: '👻', states: vec!["".to_string()] },
            Entity { id: '🚁', states: vec!["🔼💨⏳".to_string()] },
            Entity { id: '🌞', states: vec!["🌞💥⚡".to_string()] },
            Entity { id: '🦠', states: vec!["🦠🌝".to_string()] },
            Entity { id: '🏙️', states: vec!["🏙️".to_string()] },
            Entity { id: '🛰️', states: vec!["📡🔄".to_string()] }, // New entity added
        ],
        timeline: vec![],
        multiverse: vec![Universe {
            id: 0,
            entities: vec![
                Entity { id: '👻', states: vec!["".to_string()] },
                Entity { id: '🚁', states: vec!["🔼💨⏳".to_string()] },
                Entity { id: '🌞', states: vec!["🌞💥⚡".to_string()] },
                Entity { id: '🦠', states: vec!["🦠🌝".to_string()] },
                Entity { id: '🏙️', states: vec!["🏙️".to_string()] },
                Entity { id: '🛰️', states: vec!["📡🔄".to_string()] },
            ],
            observers: vec!['👻'],
        }],
    };
    
    run_simulation(&mut simulation);
    print_source();
}

fn run_simulation(sim: &mut Simulation) {
    for universe in &mut sim.multiverse {
        for i in 0..universe.entities.len() {
            if i > 0 && is_collision(&universe.entities[i - 1], &universe.entities[i]) {
                println!("⛔ Timeline disturbance detected in universe {}! Collision between {} and {}!", 
                         universe.id, universe.entities[i - 1].id, universe.entities[i].id);
                println!("🔍 Causal agent identified: {}", universe.entities[i - 1].id);
                branch_universe(sim, universe.id, universe.entities[i - 1].id);
            }
        }
    }
}

fn is_collision(entity1: &Entity, entity2: &Entity) -> bool {
    let non_matter_entities = vec!['👻']; // Define non-material entities
    let entity2_is_matter = !non_matter_entities.contains(&entity2.id); // All else is matter
    entity1.id == '🚁' && entity2_is_matter // Helicopter collides with matter
}

fn branch_universe(sim: &mut Simulation, parent_id: usize, cause: char) {
    let new_id = sim.multiverse.len();
    let new_universe = Universe {
        id: new_id,
        entities: sim.multiverse[parent_id].entities.clone(),
        observers: vec![cause], // The cause becomes an observer
    };
    println!("🌌 Branching new universe {} due to {}", new_id, cause);
    sim.multiverse.push(new_universe);
}

fn print_source() {
    let source = std::fs::read_to_string(file!()).expect("Failed to read source code");
    println!("\nQuine Output:\n\n{}", source);
}
