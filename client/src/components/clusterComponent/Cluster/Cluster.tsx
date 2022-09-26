import 'react-circular-progressbar/dist/styles.css';
import './cluster.css';
import { Row, Col, Card, Container } from 'react-bootstrap';
import ClusterHeading from '../ClusterHeader';
import ClusterResourcesChart from '../ClusterResourceChart';
import { ForeignCluster, ResourcesType } from '../../../api/types';
import { useMediaQuery } from 'react-responsive';
import React from 'react';

export interface ICluster {
  cluster: ForeignCluster;
}

function Cluster(props: ICluster) {
  const { cluster } = props;
  const isSingleCol = useMediaQuery({ query: '(max-width: 992px)' });
  const ref = React.useRef<HTMLDivElement>(null);

  return (
    <div ref={ref}>
      <Card className="margin-top">
        <Card.Header>
          <ClusterHeading cluster={cluster} />
        </Card.Header>
        <Card.Body className="p-0 m-0">
          <Container>
            <Row>
              {(!isSingleCol || cluster.outgoingResources) && (
                <Col lg={6} sm={12} className={'pt-4 pb-4'}>
                  <ClusterResourcesChart
                    resourcesType={ResourcesType.Outgoing}
                    resourcesMetrics={cluster.outgoingResources}
                    borderColor={'rgb(40, 173, 0)'}
                  />
                </Col>
              )}

              {(!isSingleCol || cluster.incomingResources) && (
                <Col lg={6} sm={12} className={'pt-4 pb-4'}>
                  <ClusterResourcesChart
                    resourcesType={ResourcesType.Incoming}
                    resourcesMetrics={cluster.incomingResources}
                    borderColor={'rgb(255, 161, 0)'}
                  />
                </Col>
              )}
            </Row>
          </Container>
        </Card.Body>
      </Card>
    </div>
  );
}

export default Cluster;
