interface Node {
  id: string;
  position: [number, number, number];
}

interface Edge {
  source: string;
  target: string;
}

export interface Graph {
  nodes: Node[];
  edges: Edge[];
}

export const sampleGraph: Graph = {
  nodes: [
    { id: '1', position: [0, 0, 0] },
    { id: '2', position: [1, 1, 1] },
    { id: '3', position: [-1, -1, -1] },
  ],
  edges: [
    { source: '1', target: '2' },
    { source: '2', target: '3' },
    { source: '3', target: '1' },
  ],
};