import React, {Component} from 'react';
import PropTypes from 'prop-types';

class Selector extends Component {

  constructor(props) {
    super(props);
    this.handleValueSelection = this.handleValueSelection.bind(this);
  }

  handleValueSelection(event) {
    this.props.selectValue(event.target.value);
  }

  render() {

    return(
      <select value={this.props.selected} onChange={this.handleValueSelection}>
        {Object.keys(this.props.values).map((value, index) => (<option key={index} value={value}>{this.props.values[value]}</option>))}
      </select>
    );
  }

}

Selector.propTypes = {
  values: PropTypes.object.isRequired,
  selected: PropTypes.oneOfType([PropTypes.string, PropTypes.number]).isRequired,
  selectValue: PropTypes.func.isRequired
};

export default Selector;
