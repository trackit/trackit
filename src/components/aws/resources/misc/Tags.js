import {Component} from "react";
import Popover from "@material-ui/core/Popover/Popover";
import PropTypes from "prop-types";
import React from "react";
import Misc from "../../../misc";

const Tooltip = Misc.Popover;

class Tags extends Component {

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
            className="tags-list"
            onClick={this.handlePopoverClose}
          >
            {Object.keys(this.props.tags).map((tag, index) => (<div key={index} className="tags-item">{tag} : {this.props.tags[tag]}</div>))}
          </div>
        </Popover>
        <div
          ref={node => {
            this.anchor = node;
          }}
          onClick={this.handlePopoverOpen}
        >
          <Tooltip placement="left" icon={<i className="fa fa-tags"/>} tooltip="Click to show tags"/>
        </div>
      </div>
    );
  }

}

Tags.propTypes = {
  tags: PropTypes.object.isRequired
};

export default Tags;
