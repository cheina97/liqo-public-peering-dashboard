export type ForeignCluster = {
  name: string;
  networking: string;
  authentication: string;
  outgoingPeering: string;
  incomingPeering: string;
  age: string;
  outgoingResources: ResourcesMetrics;
  incomingResources: ResourcesMetrics;
};

export type ResourcesMetrics = {
  totalCpus: number;
  totalMemory: number;
  usedCpus: number;
  usedMemory: number;
};

export enum ResourcesType {
  Incoming = 'Imported Resources',
  Outgoing = 'Exported Resources',
}
