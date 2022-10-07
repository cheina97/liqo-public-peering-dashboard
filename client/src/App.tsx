import 'bootstrap/dist/css/bootstrap.min.css';
import { useCallback, useEffect, useRef, useState } from 'react';
import { Col, Container, Row } from 'react-bootstrap';
import { ForeignCluster } from './api/types';
import API from './api/API';
import './App.css';
import LiqoNavbar from './components/Navbar/Navbar';
import Sidebar from './components/Sidebar/Sidebar';
import ClusterList from './components/clusterComponent/ClusterList';

function App() {
  const [clusters, setClusters] = useState<Array<ForeignCluster>>([]);
  const [currentCluster, setCurrentCluster] = useState<ForeignCluster>();
  const [init, setInit] = useState<Boolean>(true);
  const [isHamburgerOpened, setHamburgerStatus] = useState<Boolean>(false);
  const refs = useRef<Array<HTMLDivElement | null>>([]);

  const fetchAndSetCluster = useCallback(() => {
    API.getPeerings().then(clusters => {
      clusters.sort((a, b) => a.name.localeCompare(b.name));
      setClusters(clusters);
      if (clusters.length > 0) {
        if (!currentCluster) {
          setCurrentCluster(clusters[0]);
        }
      }
    });
  }, [currentCluster]);

  useEffect(() => {
    if (init) {
      fetchAndSetCluster();
      setInit(false);
    }
    const interval = setInterval(() => {
      setInit(true);
    }, 60000);
    return () => clearInterval(interval);
  }, [fetchAndSetCluster, init]);

  useEffect(() => {
    refs.current = refs.current.slice(0, clusters.length);
  }, [clusters.length]);

  useEffect(() => {
    const onScroll = () => {
      const index: number = refs.current.findIndex(ref => {
        const top: number = ref?.getBoundingClientRect().top || 0;
        return top < 300 && top > -300;
      });
      if (index !== -1 && clusters[index] !== currentCluster) {
        setCurrentCluster(clusters[index]);
      }
    };
    window.addEventListener('scroll', onScroll);
  }, [currentCluster, clusters]);

  function onClusterClick(clusterName: string) {
    const clusterIndex = clusters.findIndex(
      (cluster: ForeignCluster) => cluster.name === clusterName
    );
    if (clusterIndex !== -1) {
      refs.current[clusterIndex]?.scrollIntoView({
        behavior: 'smooth',
      });
    }
  }

  return (
    <>
      <LiqoNavbar
        onHamburgerClick={() =>
          setHamburgerStatus((oldStatus: Boolean) => !oldStatus)
        }
        isHamburgerOpened={isHamburgerOpened}
      />
      <Container fluid={true} className="navbar-padding">
        <Row>
          <Col md={2}>
            <Sidebar
              onClusterClick={onClusterClick}
              currentClusterName={currentCluster?.name}
              clustersNames={clusters.map(
                (cluster: ForeignCluster) => cluster.name
              )}
              collapsed={!isHamburgerOpened}
            />
          </Col>
          <Col md={10} className="pb-4">
            <ClusterList clusters={clusters} refs={refs} />
          </Col>
        </Row>
      </Container>
    </>
  );
}

export default App;
