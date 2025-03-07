struct Simulation {
    entities: Vec<Entity>,
    timeline: Vec<Event>,
}

struct Entity {
    id: char,
    state: String,
}

struct Event {
    timestamp: u64,
    entity_id: char,
    action: String,
}

fn main() {
    let mut simulation = Simulation {
        entities: vec![
            Entity { id: '👻', state: "" },
            Entity { id: '🚁', state: "👨‍💼🪖🔼💨⏳".to_string() },
            Entity { id: '🌞', state: "🌞💥⚡".to_string() },
            Entity { id: '🦠', state: "🦠🌝".to_string() },
            Entity { id: '🏙️', state: "🏙️👀".to_string() },
        ],
        timeline: vec![],
    };
    
    run_simulation(&mut simulation);
}

fn run_simulation(sim: &mut Simulation) {
    for i in 0..sim.entities.len() {
        if i > 0 && is_collision(&sim.entities[i - 1], &sim.entities[i]) {
            println!("⛔ Timeline disturbance detected! Collision between {} and {}!", 
                     sim.entities[i - 1].id, sim.entities[i].id);
            println!("🔍 Causal agent identified: {}", sim.entities[i - 1].id);
        }
    }
}

fn is_collision(entity1: &Entity, entity2: &Entity) -> bool {
    let non_matter_entities = vec!['👻']; // Define non-material entities
    let entity2_is_matter = !non_matter_entities.contains(&entity2.id); // All else is matter
    entity1.id == '🚁' && entity2_is_matter // Helicopter collides with matter
}