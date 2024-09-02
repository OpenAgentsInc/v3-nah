import React from 'react';
import { View, StyleSheet } from 'react-native';
import { Canvas } from '@react-three/fiber/native';
import { OrbitControls } from '@react-three/drei/native';

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
    <mesh position={position}>
      <sphereGeometry args={[0.1, 32, 32]} />
      <meshStandardMaterial color="white" />
    </mesh>
  );
};

const Edge: React.FC<{ start: [number, number, number]; end: [number, number, number] }> = ({ start, end }) => {
  return (
    <line>
      <bufferGeometry attach="geometry">
        <bufferAttribute
          attachObject={['attributes', 'position']}
          count={2}
          array={new Float32Array([...start, ...end])}
          itemSize={3}
        />
      </bufferGeometry>
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