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
    return(
      <Selector values={intervals} selected={this.props.interval} selectValue={this.props.setInterval}/>
    );
  }

}

IntervalSelector.propTypes = {
  interval: PropTypes.string,
  setInterval: PropTypes.func
};

export default IntervalSelector;
