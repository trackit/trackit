import {Component} from "react";
import Popover from "@material-ui/core/Popover/Popover";
import {formatBytes} from "../../../../common/formatters";
import PropTypes from "prop-types";
import React from "react";
import Misc from "../../../misc";

const Tooltip = Misc.Popover;

class UnusedStorage extends Component {

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
    return (
      <div>
        <Popover
          open={this.state.showPopOver}
          anchorEl={this.anchor}
          onClose={this.handlePopoverClose}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'center',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'right',
          }}
        >
          <div
            className="unusedStorage-list"
            onClick={this.handlePopoverClose}
          >
            {Object.keys(this.props.data).map((item, index) => (<div key={index} className="unusedStorage-item">{item} : {this.props.data[item] >= 0 ? formatBytes(this.props.data[item]) : "N/A"}</div>))}
          </div>
        </Popover>
        <div
          ref={node => {
            this.anchor = node;
          }}
          onClick={this.handlePopoverOpen}
        >
          <Tooltip placement="right" info tooltip="Click to see more details"/>
        </div>
      </div>
    );
  }

}

UnusedStorage.propTypes = {
  data: PropTypes.object.isRequired
};

export default UnusedStorage;
