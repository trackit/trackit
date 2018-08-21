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
    ) : this.props.icon);

    let child = null;
    if (this.props.children)
      child = this.props.children;
    else if (this.props.tooltip)
      child = (<Tooltip className="popover-tooltip" id="tooltip">{this.props.tooltip}</Tooltip>);

    return (
      <div className="wrapper">
        <OverlayTrigger placement={this.props.placement} overlay={child}>
          <div className="popover-trigger">{trigger}</div>
        </OverlayTrigger>
      </div>
    )
  }
}

Popover.propTypes = {
  icon: PropTypes.node,
  info: PropTypes.bool,
  children: PropTypes.node,
  tooltip: PropTypes.node,
  placement: PropTypes.oneOf(["top", "bottom", "left", "right"]),
  triggerStyle: PropTypes.object,
};

Popover.defaultProps = {
  info: false,
  children: null,
  tooltip: null,
  placement: "right"
};

export default Popover;
