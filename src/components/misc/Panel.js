import React, {Component} from 'react';
import PropTypes from "prop-types";

export class PanelItem extends Component {

  render() {

    const classes = "white-box " + (this.props.children && this.props.children.props && this.props.children.props.className ? this.props.children.props.className : "");

    return (
      <div className="row">
        <div className="col-md-12">
          <div className={classes}>
            {this.props.children}
          </div>
        </div>
      </div>
    );
  }

}

PanelItem.propTypes = {
  children: PropTypes.node.isRequired
};

class Panel extends Component {

  render() {
    let body = ((Array.isArray(this.props.children)) ?
      this.props.children.map((item, index) => ((item !== null ? <PanelItem key={index} children={item}/> : null))) :
      <PanelItem children={this.props.children}/>
    );
    return(
      <div className="container-fluid">
        {body}
      </div>
    );
  }

}

Panel.propTypes = {
  children: PropTypes.node.isRequired
};


export default Panel;
