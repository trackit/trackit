import React, { Component } from 'react';
import PropTypes from 'prop-types';

// Chart Component
class BarChartComponent extends Component {

    componentDidMount() {
      var layout = this.getLayout(this.props.barmode);

      window.Plotly.newPlot(this.props.elementId, this.props.data, layout, {displaylogo: false});
    }

    componentWillReceiveProps(nextProps) {
      var layout = this.getLayout(nextProps.barmode);
      window.Plotly.purge(nextProps.elementId);
      window.Plotly.plot(nextProps.elementId, nextProps.data, layout, {displaylogo: false});
    }

    getLayout(barmode) {
      return(
        {
          barmode,
          xaxis: {title: 'Regions'},
          yaxis: {title: 'Price ($)'},
          height: 350,
          title: this.props.title,
        }
      );
    }

    render() {


      return (
        <div>
          <div id={this.props.elementId} />
        </div>
      );
    }

}

BarChartComponent.defaultProps = {
  barmode: 'group',
  title: '',
}

// Define PropTypes
BarChartComponent.propTypes = {
  elementId: PropTypes.string.isRequired,
  data: PropTypes.array.isRequired,
  barmode: PropTypes.string,
  title: PropTypes.string
};


export default BarChartComponent;
