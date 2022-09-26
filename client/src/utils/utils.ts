import dayjs from 'dayjs';
import duration, { Duration } from 'dayjs/plugin/duration';
import { ResourcesType } from '../api/types';

dayjs.extend(duration);

const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

export function getDuration(clusterAge: string): string {
  if (!clusterAge) {
    return 'None';
  }
  const duration: Duration = getDurationFromDate(clusterAge);
  const result = [];
  const days = duration.days();
  const hours = duration.hours();
  const minutes = duration.minutes();

  if (days > 0) result.push(`${days} days`);
  if (hours > 0) result.push(` ${hours} hours`);
  if (days === 0 || hours === 0) result.push(` ${minutes} minutes`);

  return result.join(' ');
}

export function getDurationFromDate(dateString: string): Duration {
  const peeringCreationDate = getIsoFormattedDate(dateString);
  return dayjs.duration(
    dayjs().diff(dayjs(peeringCreationDate, 'YYYY-MM-DD HH:mm:ss Z'))
  );
}

export function calculatePercentage(total: number, actual: number): number {
  const percentage = (actual / total) * 100;
  return percentage > 100 ? 100 : percentage;
}

export function getHighestUnit(bytes: number): string {
  let unitIndex = 0;
  while (bytes >= 1024) {
    bytes /= 1024;
    unitIndex++;
  }
  return `${bytes?.toFixed(1)} ${units[unitIndex]}`;
}

export function textOnChart(
  chart: any,
  text: string,
  scaleFactor: number,
  bold: boolean
): void {
  var fontSize = (chart.height / 350).toFixed(2);
  chart.ctx.font = `${bold ? 'bold ' : ''} ` + fontSize + 'em sans-serif';
  chart.ctx.textBaseline = 'top';
  var textX = Math.round((chart.width - chart.ctx.measureText(text).width) / 2);
  var textY = chart.height / scaleFactor;
  chart.ctx.fillText(text, textX, textY);
}

export function noResourcesMessage(type: ResourcesType): string {
  if (type === ResourcesType.Incoming) {
    return 'No Imported Resources (incoming peering disabled)';
  }
  return 'No Exported Resources (outgoing peering disabled)';
}

function removeCharAt(index: number, s: string): string {
  return s.substring(0, index) + s.substring(index + 1, s.length);
}

function getIsoFormattedDate(s: string): string {
  const dateWithoutTimezone: string = s.slice(0, -4);
  const lastSpaceIndex: number = dateWithoutTimezone.lastIndexOf(' ');
  return removeCharAt(lastSpaceIndex, dateWithoutTimezone);
}
