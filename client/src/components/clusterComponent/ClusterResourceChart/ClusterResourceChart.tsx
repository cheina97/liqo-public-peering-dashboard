import { Doughnut } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  TooltipItem,
} from 'chart.js';
import { Col, Container, Row } from 'react-bootstrap';
import {
  calculatePercentage,
  getHighestUnit,
  noResourcesMessage,
  textOnChart,
} from '../../../utils/utils';
import { ResourcesMetrics, ResourcesType } from '../../../api/types';

ChartJS.register(ArcElement, Tooltip, Legend);

export interface IClusterResourcesChart {
  resourcesType: ResourcesType;
  resourcesMetrics: ResourcesMetrics;
  borderColor: string;
}

const backgroundColors = [
  'rgba(54, 162, 235, 0.5)',
  'rgba(255, 162, 235, 0.5)',
  'rgba(255, 255, 255, 1)',
];

function ClusterResourcesChart(props: IClusterResourcesChart) {
  const { resourcesType, resourcesMetrics, borderColor } = props;

  const cpuUsage = calculatePercentage(
    resourcesMetrics?.totalCpus,
    resourcesMetrics?.usedCpus
  );

  const memUsage = calculatePercentage(
    resourcesMetrics?.totalMemory,
    resourcesMetrics?.usedMemory
  );

  const data = {
    labels: ['CPU', 'Memory', 'Free'],
    datasets: [
      {
        label: 'CPU',
        data: [cpuUsage, 0, 100 - cpuUsage],
        backgroundColor: backgroundColors,
        borderColor: borderColor,
        borderWidth: 0.6,
      },
      {
        label: 'Memory',
        data: [0, memUsage, 100 - memUsage],
        backgroundColor: backgroundColors,
        borderColor: borderColor,
        borderWidth: 0.6,
      },
    ],
  };

  const cpuUsageString: string =
    resourcesMetrics?.usedCpus?.toFixed(1) +
    ' / ' +
    resourcesMetrics?.totalCpus?.toFixed(1);

  const memUsageString: string =
    getHighestUnit(resourcesMetrics?.usedMemory) +
    ' / ' +
    getHighestUnit(resourcesMetrics?.totalMemory);
  var plugins = [
    {
      id: '1',
      beforeDraw: function (chart: any) {
        textOnChart(chart, 'vCPUs used:', 2.0, true);
        textOnChart(chart, cpuUsageString, 1.8, false);
        textOnChart(chart, 'Memory used:', 1.5, true);
        textOnChart(chart, memUsageString, 1.38, false);
      },
    },
    {
      id: '2',
      beforeInit: function (chart: any) {
        const originalFit = chart.legend.fit;
        chart.legend.fit = function fit() {
          originalFit.bind(chart.legend)();
          this.height += this.height / 1.5;
        };
      },
    },
  ];

  return resourcesMetrics ? (
    <>
      <Row className="justify-content-center textshadow text h4 mb-4">
        {resourcesType}
      </Row>
      <Row>
        <Col>
          <Container style={{ minHeight: '320px' }}>
            <Doughnut
              data={data}
              options={{
                animation: false,
                circumference: 270,
                rotation: -135,
                responsive: true,
                maintainAspectRatio: false,
                cutout: '70%',
                plugins: {
                  legend: {
                    labels: {
                      usePointStyle: true,
                      pointStyle: 'rect',
                      filter(item, _) {
                        return item.text !== 'Free';
                      },
                    },
                  },
                  tooltip: {
                    callbacks: {
                      label: (data: TooltipItem<'doughnut'>) => {
                        return `${data.label}: ${data.formattedValue}%`;
                      },
                    },
                  },
                },
              }}
              plugins={plugins}
            />
          </Container>
        </Col>
      </Row>
    </>
  ) : (
    <div className="text-center h5 vcenter h-100">
      {noResourcesMessage(resourcesType)}
    </div>
  );
}

export default ClusterResourcesChart;
