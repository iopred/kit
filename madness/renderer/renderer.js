
// Create a scene
const scene = new THREE.Scene();

// Create a camera
const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
camera.position.z = 5;

// Create a renderer
const renderer = new THREE.WebGLRenderer();
renderer.setSize(window.innerWidth, window.innerHeight);
document.body.appendChild(renderer.domElement);

// Create a cube and add it to the scene
const geometry = new THREE.BoxGeometry( 1, 1, 1 );
const material = new THREE.MeshPhysicalMaterial({ color: 0xff1111 });
const cube = new THREE.Mesh(geometry, material);

// Enable shadow casting and receiving
cube.castShadow = true;
cube.receiveShadow = true;

scene.add(cube);

const directionalLight = new THREE.DirectionalLight(0xffffff, 1);
directionalLight.position.set(5, 5, 5);
directionalLight.castShadow = true;
scene.add(directionalLight);

// Create a floor geometry
const floorGeometry = new THREE.PlaneGeometry(10, 10);

// Generate a checkerboard texture procedurally
const checkerboardTexture = new THREE.DataTexture(new Uint8Array([255, 255, 255, 255]), 1, 1);
checkerboardTexture.wrapS = THREE.RepeatWrapping;
checkerboardTexture.wrapT = THREE.RepeatWrapping;
checkerboardTexture.repeat.set(4, 4);

// Create a material with the checkerboard texture
const floorMaterial = new THREE.MeshPhysicalMaterial({color: 0xffff00, side: THREE.DoubleSide });

// Create the floor mesh
const floor = new THREE.Mesh(floorGeometry, floorMaterial);
floor.rotation.x = -Math.PI / 2; // Rotate the floor to be horizontal
floor.position.y -= 1;
floor.receiveShadow = true; // Make the floor receive shadows

// Add the floor to the scene
scene.add(floor);


// Enable shadow mapping
renderer.shadowMap.enabled = true;
renderer.shadowMap.type = THREE.PCFSoftShadowMap;
scene.background = new THREE.Color("rgb(32, 32, 32)");

// Create an animation function
const animate = () => {
  requestAnimationFrame(animate);

  // Rotate the cube
  cube.rotation.x += 0.01;
  cube.rotation.y += 0.01;

  // Render the scene
  renderer.render(scene, camera);
};

// Start the animation
animate();	




