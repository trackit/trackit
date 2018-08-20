import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { Tooltip, OverlayTrigger } from 'react-bootstrap';
import '../../styles/Popover.css';

class Popover extends Component {

  constructor(props) {
    super(props);
    this.state = {
      showPopOver: false
    };
    this.handlePopoverOpen = this.handlePopoverOpen.bind(this);
    this.handlePopoverClose = this.handlePopoverClose.bind(this);
  }

  handlePopoverOpen = (e) => {
    e.preventDefault();
    this.setState({ showPopOver: true });
  };

  handlePopoverClose = (e) => {
    e.preventDefault();
    this.setState({ showPopOver: false });
  };

  render() {
    const trigger = (this.props.info ? (
      <i className="fa fa-info-circle" style={this.props.triggerStyle}/>
    ) : this.props.children);

    const tooltip = (<Tooltip className="popover-tooltip" id="tooltip">{this.props.popOver}</Tooltip>);

    return (
      <div className="wrapper">
        <OverlayTrigger placement="right" overlay={tooltip}>
          <div className="popover-trigger">{trigger}</div>
        </OverlayTrigger>
      </div>
    )
  }
}

Popover.propTypes = {
  children: PropTypes.node,
  info: PropTypes.bool,
  popOver: PropTypes.node.isRequired,
  triggerStyle: PropTypes.object,
};

Popover.defaultProps = {
  info: false
};

export default Popover;
