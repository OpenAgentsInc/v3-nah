import React from 'react';
import { View, StyleSheet } from 'react-native';
import { Canvas, useThree } from '@react-three/fiber/native';
import { OrbitControls } from '@react-three/drei/native';
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

const Node: React.FC<{ position: [number, number, number] }> = ({ position }) => {
  return (
    <group position={position}>
      <mesh>
        <sphereGeometry args={[0.1, 32, 32]} />
        <meshBasicMaterial color="black" />
      </mesh>
      <lineSegments>
        <sphereGeometry args={[0.1, 32, 32]} />
        <lineBasicMaterial color="white" />
      </lineSegments>
    </group>
  );
};

const Edge: React.FC<{ start: [number, number, number]; end: [number, number, number] }> = ({ start, end }) => {
  const points = [new THREE.Vector3(...start), new THREE.Vector3(...end)];
  const lineGeometry = new THREE.BufferGeometry().setFromPoints(points);

  return (
    <line geometry={lineGeometry}>
      <lineBasicMaterial attach="material" color="white" linewidth={1} />
    </line>
  );
};

const Graph: React.FC<GraphCanvasProps> = ({ nodes, edges }) => {
  return (
    <>
      {nodes.map((node) => (
        <Node key={node.id} position={node.position} />
      ))}
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
    </>
  );
};

const GraphCanvas: React.FC<GraphCanvasProps> = ({ nodes, edges }) => {
  return (
    <View style={styles.canvasContainer}>
      <Canvas style={styles.canvas}>
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
    zIndex: -1,
  },
  canvas: {
    flex: 1,
  },
});

export default GraphCanvas;