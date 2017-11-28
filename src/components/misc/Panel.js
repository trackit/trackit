import React, {Component} from 'react';
import PropTypes from "prop-types";

class PanelItem extends Component {

  render() {
    return (
      <div className="row">
        <div className="col-md-12">
          <div className="white-box">
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
      this.props.children.map((item, index) => (<PanelItem key={index} children={item}/>)) :
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
