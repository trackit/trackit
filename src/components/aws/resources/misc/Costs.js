import {Component} from "react";
import Popover from "@material-ui/core/Popover/Popover";
import PropTypes from "prop-types";
import React from "react";
import {formatPrice} from "../../../../common/formatters";
import Misc from "../../../misc";

const Tooltip = Misc.Popover;

class Costs extends Component {

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
            className="costs-list"
            onClick={this.handlePopoverClose}
          >
            {Object.keys(this.props.costs).map((tag, index) => (<div key={index} className="costs-item">{tag} : {formatPrice(this.props.costs[tag])}</div>))}
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

Costs.propTypes = {
  costs: PropTypes.object.isRequired
};

export default Costs;
