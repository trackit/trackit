import React, {Component} from 'react';
import moment from 'moment';
import PropTypes from 'prop-types';
import NVD3Chart from 'react-nvd3';
import * as d3 from 'd3';

const context = {
    formatXAxis: (d) => (d3.time.format('%x')(new Date(d))),
    formatYAxis: (d) => ('$' + d3.format(',.2f')(d)),
  };
  
  const xAxis = {
    tickFormat: {
      name:'formatXAxis',
      type:'function',
    }
  };
  
  const yAxis = {
    tickFormat: {
      name:'formatYAxis',
      type:'function',
    }
  };
  
  /* istanbul ignore next */
  const formatX = (d) => {
    const date = new Date(d[0]);
    return date.getTime();
  };
  
  /* istanbul ignore next */
  const formatY = (d) => (d[1]);
  
  const margin = {
    right: 10,
    left: 70,
  };
  

class EventPanel extends Component {
    formatDataForChart(data, service) {
        const res = [
          {
            key: `${service.length ? service : "Unknown service"} cost`,
            values: [],
          },
          {
            key: "Anomaly",
            values: [],
            color: '#ff0000'
          },
        ];
        for (let i = 0; i < data.length; i++) {
          const element = data[i];
          res[0].values.push([element.date, element.cost]);
          res[1].values.push([element.date, element.abnormal ? element.cost : 0]);
        }
        return res;
    }

    isolateAnomaly(data, anomaly) {
        const res = JSON.parse(JSON.stringify(data));
        for (let i = 0; i < res.length; i++) {
            const element = res[i];
            if (!(element.abnormal && element.date === anomaly.date)) {
                res[i].abnormal = false;
            }
        }
        return res;
    }

    getBadgeClasses(level) {
        switch (level) {
            case 0:
                return 'badge green-bg'
            case 1:
                return 'badge orange-bg'
            case 2:
                return 'badge red-bg'
            case 3:
                return 'badge red-bg'
            default:
                return 'badge red-bg'
        }
    }
    

    render() {
        const { dataSet, abnormalElement, service } = this.props;
        const exceededCost = (abnormalElement.cost - abnormalElement.upper_band).toFixed(2);
        const badgeLabels = ['Low', 'Medium', 'High', 'Critical'];
        const anomalyLevel = abnormalElement.level;
        return (
            <div className="white-box">
                <h5 className="inline-block">
                    <i className="fa fa-exclamation-circle"></i>
                    &nbsp;
                    {service.length ? service : "Unknown service"}
                    &nbsp;
                    <span className={this.getBadgeClasses(anomalyLevel)}>{badgeLabels[anomalyLevel]}</span>
                </h5>
                <h5 className="inline-block pull-right">{moment(abnormalElement.date).format("ddd, MMM Do Y")}</h5>
                <div className="clearfix"></div>
                <p>On {moment(abnormalElement.date).format("ddd, MMM Do Y")}, <strong>{service.length ? service : "Unknown service"}</strong> exceeded the maximum expected cost for this service by <strong>${exceededCost}</strong></p>
                <hr />
                <NVD3Chart
                    id="barChart"
                    type="multiBarChart"
                    datum={this.formatDataForChart(this.isolateAnomaly(dataSet, abnormalElement), service)}
                    context={context}
                    xAxis={xAxis}
                    yAxis={yAxis}
                    margin={margin}
                    rightAlignYAxis={false}
                    clipEdge={false}
                    showControls={false}
                    showLegend={true}
                    stacked={false}
                    x={formatX}
                    y={formatY}
                    height={250}
                />
            </div>
        );
    }
}

EventPanel.propTypes = {
    dataSet: PropTypes.array.isRequired,
    abnormalElement: PropTypes.object.isRequired,
    service: PropTypes.string.isRequired,
};

export default EventPanel;