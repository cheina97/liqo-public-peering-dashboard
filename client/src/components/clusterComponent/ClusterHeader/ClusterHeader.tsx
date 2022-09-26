import { Col, Container, Row } from 'react-bootstrap';
import { getDuration } from '../../../utils/utils';
import { useMediaQuery } from 'react-responsive';
import { ForeignCluster } from '../../../api/types';

export interface IClusterHeader {
  cluster: ForeignCluster;
}

function ClusterHeading(props: IClusterHeader) {
  const { cluster } = props;
  const isTwoLines = useMediaQuery({ query: '(max-width: 960px)' });
  const isMobile = useMediaQuery({ query: '(max-width: 768px)' });

  return (
    <Container>
      <Col>
        <Row className="text-center vcenter pb-2 mb-2 mt-1 h3 border-bottom">
          {cluster.name}
        </Row>
        <Row className="justify-content-center">
          <Col ms={12} md={4} lg={2}>
            <Row>
              <Col ms={6} md={12} className="vcenter textshadow">
                Networking
              </Col>
              <Col ms={6} md={12} className="vcenter">
                {cluster.networking}
              </Col>
            </Row>
          </Col>
          <Col ms={12} md={4} lg={2}>
            <Row>
              <Col ms={6} md={12} className="vcenter textshadow">
                Authentication
              </Col>
              <Col ms={6} md={12} className="vcenter">
                {cluster.authentication}
              </Col>
            </Row>
          </Col>
          <Col ms={12} md={4} lg={2}>
            <Row>
              <Col ms={6} md={12} className="vcenter textshadow">
                Outgoing Peering
              </Col>
              <Col ms={6} md={12} className="vcenter">
                {cluster.outgoingPeering}
              </Col>
            </Row>
          </Col>
          <Col
            ms={12}
            md={4}
            lg={2}
            className={isTwoLines && !isMobile ? 'mt-3' : ''}
          >
            <Row>
              <Col ms={6} md={12} className="vcenter textshadow">
                Incoming Peering
              </Col>
              <Col ms={6} md={12} className="vcenter">
                {cluster.incomingPeering}
              </Col>
            </Row>
          </Col>
          <Col
            ms={12}
            md={4}
            lg={2}
            className={isTwoLines && !isMobile ? 'mt-3' : ''}
          >
            <Row>
              <Col ms={6} md={12} className="vcenter textshadow">
                Up Time
              </Col>
              <Col ms={6} md={12} className="vcenter">
                {getDuration(cluster.age)}
              </Col>
            </Row>
          </Col>
        </Row>
      </Col>
    </Container>
  );
}

export default ClusterHeading;
