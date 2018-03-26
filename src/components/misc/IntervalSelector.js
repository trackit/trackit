import React, {Component} from 'react';
import PropTypes from 'prop-types';

import Selector from './Selector';

const intervals = {
  day: "Daily",
  week: "Weekly",
  month: "Monthly",
  year: "Yearly"
};

class IntervalSelector extends Component {

  render() {

    const listedIntervals = (!this.props.availableIntervals || !this.props.availableIntervals.length ? intervals : {});
    if (this.props.availableIntervals) {
      this.props.availableIntervals.forEach((interval) => {
        listedIntervals[interval] = intervals[interval];
      });
    }

    return(
      <Selector values={listedIntervals} selected={this.props.interval} selectValue={this.props.setInterval}/>
    );
  }

}

IntervalSelector.propTypes = {
  interval: PropTypes.string,
  setInterval: PropTypes.func,
  availableIntervals: PropTypes.arrayOf(PropTypes.string)
};

export default IntervalSelector;
