import React, { useRef, useMemo } from 'react';
import { View, StyleSheet } from 'react-native';
import { Canvas, useFrame } from '@react-three/fiber/native';
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

const NODE_RADIUS = 0.1;

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

  const edgeGeometry = useMemo(() => {
    const geometry = new THREE.BufferGeometry();
    const positions = new Float32Array([
      adjustedStart.x, adjustedStart.y, adjustedStart.z,
      adjustedEnd.x, adjustedEnd.y, adjustedEnd.z
    ]);
    geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    return geometry;
  }, [adjustedStart, adjustedEnd]);

  return (
    <lineSegments geometry={edgeGeometry}>
      <lineBasicMaterial color="white" />
    </lineSegments>
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