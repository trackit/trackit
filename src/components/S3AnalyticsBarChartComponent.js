import React, { Component } from 'react';
import PropTypes from 'prop-types';

 // S3AnalyticsBarChart Component
 class S3AnalyticsBarChart extends Component {

   formatDataForChart(data) {
     // Cloning the data
     let dataClone = data.slice(0);
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
     for (let i = 0; i < dataClone.length; i += 1) {
       const tmp = dataClone[i];
       bandwidth.x.push(tmp._id);
       storage.x.push(tmp._id);
       bandwidth.y.push(tmp.bw_cost.toFixed(2));
       storage.y.push(tmp.storage_cost.toFixed(2));
     }
     return [bandwidth, storage];
   }

   componentDidMount() {
     const data = this.formatDataForChart(this.props.data);
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

S3AnalyticsBarChart.propTypes = {
  elementId: PropTypes.string.isRequired,
  data: PropTypes.array.isRequired,
};

 export default S3AnalyticsBarChart;
