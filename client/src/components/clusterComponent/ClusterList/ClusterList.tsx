import { ForeignCluster } from '../../../api/types';
import Cluster from '../Cluster/Cluster';

export interface IClusterList {
  clusters: ForeignCluster[];
  refs: React.MutableRefObject<(HTMLDivElement | null)[]>;
}

export default function ClusterList(props: IClusterList) {
  const { clusters, refs } = props;

  return (
    <>
      {clusters.map((cluster: ForeignCluster, i: number) => (
        <div key={i} ref={el => (refs.current[i] = el)}>
          <Cluster cluster={cluster} />
        </div>
      ))}
    </>
  );
}
