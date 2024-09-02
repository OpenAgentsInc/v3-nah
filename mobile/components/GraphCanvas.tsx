import React, { useRef, useMemo } from 'react';
import { View, StyleSheet } from 'react-native';
import { Canvas, useFrame, extend } from '@react-three/fiber/native';
import { OrbitControls, shaderMaterial } from '@react-three/drei/native';
import * as THREE from 'three';

interface Node {
  id: string;
  position: [number, number, number];
}

interface Edge {
  source: string;
  target: string;
}

interface GraphCanvasProps {
  nodes: Node[];
  edges: Edge[];
}

const NODE_RADIUS = 0.1;

const LineShaderMaterial = shaderMaterial(
  { color: new THREE.Color('white') },
  // Vertex Shader
  `
    void main() {
      gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
    }
  `,
  // Fragment Shader
  `
    uniform vec3 color;
    void main() {
      gl_FragColor = vec4(color, 1.0);
    }
  `
);

extend({ LineShaderMaterial });

const Node: React.FC<{ position: [number, number, number] }> = ({ position }) => {
  const ref = useRef<THREE.Mesh>(null);

  useFrame(({ camera }) => {
    if (ref.current) {
      ref.current.quaternion.copy(camera.quaternion);
    }
  });

  return (
    <group position={position}>
      <mesh ref={ref}>
        <circleGeometry args={[NODE_RADIUS, 32]} />
        <meshBasicMaterial color="black" />
      </mesh>
      <mesh ref={ref}>
        <ringGeometry args={[NODE_RADIUS - 0.01, NODE_RADIUS, 32]} />
        <meshBasicMaterial color="white" />
      </mesh>
    </group>
  );
};

const Edge: React.FC<{ start: [number, number, number]; end: [number, number, number] }> = ({ start, end }) => {
  const startVec = new THREE.Vector3(...start);
  const endVec = new THREE.Vector3(...end);
  const direction = endVec.clone().sub(startVec).normalize();
  
  const adjustedStart = startVec.clone().add(direction.clone().multiplyScalar(NODE_RADIUS));
  const adjustedEnd = endVec.clone().sub(direction.clone().multiplyScalar(NODE_RADIUS));

  const points = useMemo(() => [adjustedStart, adjustedEnd], [adjustedStart, adjustedEnd]);
  const geometry = useMemo(() => new THREE.BufferGeometry().setFromPoints(points), [points]);

  return (
    <line geometry={geometry}>
      {/* @ts-ignore */}
      <lineShaderMaterial attach="material" color="white" />
    </line>
  );
};

const Graph: React.FC<GraphCanvasProps> = ({ nodes, edges }) => {
  return (
    <>
      {edges.map((edge, index) => {
        const sourceNode = nodes.find((n) => n.id === edge.source);
        const targetNode = nodes.find((n) => n.id === edge.target);
        if (sourceNode && targetNode) {
          return (
            <Edge
              key={`${edge.source}-${edge.target}`}
              start={sourceNode.position}
              end={targetNode.position}
            />
          );
        }
        return null;
      })}
      {nodes.map((node) => (
        <Node key={node.id} position={node.position} />
      ))}
    </>
  );
};

const GraphCanvas: React.FC<GraphCanvasProps> = ({ nodes, edges }) => {
  return (
    <View style={styles.canvasContainer}>
      <Canvas>
        <color attach="background" args={["#000000"]} />
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} />
        <Graph nodes={nodes} edges={edges} />
        <OrbitControls enablePan={true} enableZoom={true} enableRotate={true} />
      </Canvas>
    </View>
  );
};

const styles = StyleSheet.create({
  canvasContainer: {
    ...StyleSheet.absoluteFillObject,
  },
});

export default GraphCanvas;