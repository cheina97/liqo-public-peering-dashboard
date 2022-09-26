import './sidebar.css';
import { Col, Nav, Row } from 'react-bootstrap';

export interface ISidebar {
  clustersNames: Array<string>;
  currentClusterName: string | undefined;
  onClusterClick: (clusterName: string) => void;
  collapsed: boolean;
}

function Sidebar(props: ISidebar) {
  const { clustersNames, currentClusterName, onClusterClick, collapsed } =
    props;

  const clusters = clustersNames.map((clusterName: string) => (
    <Col key={clusterName} md={12} sm={12}>
      <Nav.Item>
        <Nav.Link
          className="text-center m-2 text-white"
          eventKey={clusterName}
          onClick={() => onClusterClick(clusterName)}
        >
          {clusterName}
        </Nav.Link>
      </Nav.Item>
    </Col>
  ));

  return (
    <Nav
      className="col-md-2 col-sm-12 d-md-block bg-dark sidebar justify-content-center"
      activeKey={currentClusterName || ''}
      variant="pills"
      hidden={collapsed}
    >
      <Row className="justify-content-center">
        <Col md={12} sm={12}>
          <Nav.Item className="text-light text-center">
            <b>Available Clusters</b>
          </Nav.Item>
        </Col>
        {clusters}
      </Row>
    </Nav>
  );
}

export default Sidebar;
