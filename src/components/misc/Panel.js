import React, {Component} from 'react';
import PropTypes from "prop-types";

class Panel extends Component {

  constructor(props) {
    super(props);
    this.state = {
      collapsed: props.defaultCollapse
    };
    this.toggleCollapse = this.toggleCollapse.bind(this);
  }

  toggleCollapse = (e) => {
    e.preventDefault();
    if (this.props.collapsible) {
      const collapsed = !this.state.collapsed;
      this.setState({collapsed});
    }
  };

  render() {
    const body = ((!this.props.collapsible || !this.state.collapsed) ? (
      <div className="panel-body">
        {this.props.children}
      </div>
    ) : null);
    const collapseIcon = (this.props.collapsible ? (
      <div className="pull-right">
        <span className={"glyphicon glyphicon-chevron-" + (this.state.collapsed ? "down" : "up")} aria-hidden="true"/>
      </div>
    ) : null);
    return(
      <div className="panel panel-default">

        <div className="panel-heading" onClick={this.toggleCollapse}>
          <h3 className="panel-title pull-left">{this.props.title}</h3>
          {collapseIcon}
          <div className="clearfix"/>
        </div>

        {body}

      </div>
    );
  }

}

Panel.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  collapsible: PropTypes.bool,
  defaultCollapse: PropTypes.bool
};

Panel.defaultProps = {
  collapsible: false,
  defaultCollapse: false
};


export default Panel;
