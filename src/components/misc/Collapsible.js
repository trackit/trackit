import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Collapse from "@material-ui/core/Collapse/Collapse";
import '../../styles/Collapsible.css';

class Collapsible extends Component {

  constructor(props) {
    super(props);
    this.state = {
      expanded: false
    };
    this.toggleCollapse = this.toggleCollapse.bind(this);
  }

  toggleCollapse = (e) => {
    e.preventDefault();
    this.setState({ expanded: !this.state.expanded });
  };

  render() {
    return (
      <div className={"collapsible " + this.props.className}>
        <div className="collapsible-header">
          <div className="collapsible-header-content">
            {this.props.header}
          </div>
          {(this.props.children ? (this.state.expanded ? <ExpandLess onClick={this.toggleCollapse}/> : <ExpandMore onClick={this.toggleCollapse}/>) : null)}
        </div>
        {(this.props.children ? (
          <Collapse className="collapsible-content" in={this.state.expanded} timeout={this.props.timeout} unmountOnExit={this.props.unmountOnExit}>
            {this.props.children}
          </Collapse>
        ) : null)}
      </div>
    )
  }
}

Collapsible.propTypes = {
  className: PropTypes.string,
  header: PropTypes.node,
  children: PropTypes.node,
  timeout: PropTypes.number,
  unmountOnExit: PropTypes.bool
};

Collapsible.defaultProps = {
  header: null,
  children: null,
  timeout: 0,
  unmountOnExit: false
};

export default Collapsible;
