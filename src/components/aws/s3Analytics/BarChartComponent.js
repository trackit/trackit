import React, { Component } from 'react';
import PropTypes from 'prop-types';

 // S3AnalyticsBarChart Component
 class BarChart extends Component {

   formatDataForChart() {
     // Cloning the data
     let dataClone = this.props.data.slice(0);
     // Extracting X first buckets by price
     dataClone = dataClone.sort((a,b) => {
        if (a.total_cost < b.total_cost)
          return 1;
        return -1;
     });
     dataClone = dataClone.slice(0,25);

     // Formatting for chart
     const bandwidth = {
       x: [],
       y: [],
       name: 'Bandwidth',
       type: 'bar',
       opacity: 0.8,
       marker: {
         color: '#1e88e5 ',
       }
     };
     const storage = {
       x: [],
       y: [],
       name: 'Storage',
       type: 'bar',
       opacity: 0.8,
       marker: {
         color: '#ff9800',
       },
       hoverlabel: {
         bordercolor: '#ffffff',
       }
     };

     dataClone.forEach((item) => {
       bandwidth.x.push(item._id);
       bandwidth.y.push(item.bw_cost.toFixed(2));
       storage.x.push(item._id);
       storage.y.push(item.storage_cost.toFixed(2));
     });

     return [bandwidth, storage];
   }

   componentDidMount() {
     const data = this.formatDataForChart();
     const layout = {
       barmode: 'stack',
       showlegend: true,
       title: 'Buckets breakdown',
       height: 180,
       margin: {
         l: 55,
         r: 45,
         b: 55,
         t: 25,
       },
       autosize: true,
       yaxis: {
         title: 'Total Price ($)',
       }
     };

     window.Plotly.newPlot(this.props.elementId, data, layout, {displayModeBar: false});
   }

   render() {

     return(
       <div>
         <div id={this.props.elementId}></div>
       </div>
     );
   }
 }

BarChart.propTypes = {
  elementId: PropTypes.string.isRequired,
  data: PropTypes.arrayOf(
    PropTypes.shape({
      _id: PropTypes.string.isRequired,
      size: PropTypes.number.isRequired,
      storage_cost: PropTypes.number.isRequired,
      bw_cost: PropTypes.number.isRequired,
      total_cost: PropTypes.number.isRequired,
      transfer_in: PropTypes.number.isRequired,
      transfer_out: PropTypes.number.isRequired,
    })
  ),
};

export default BarChart;
