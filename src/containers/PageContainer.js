import React, { Component } from 'react';
import { connect } from 'react-redux';
// import PropTypes from 'prop-types';

import Actions from '../actions';

import Paper from 'material-ui/Paper';

import Components from '../components';

const DeclareSetup = Components.DeclareSetup;

const style = {
  paper: {
    padding: '15px',
    margin: '15px',
  },
  noPadding: {
    padding: '0px'
  }
};

// PageContainer Component
class PageContainer extends Component {

  componentDidMount() {
    this.props.getGCPPricing();
    this.props.getAWSPricing();
  }

  render() {
    return (
      <div className="container-fluid">
        <Paper elevation={3} style={style.paper} className="animated bounceInRight">
          <DeclareSetup />
        </Paper>
      </div>
    );
  }

}

PageContainer.propTypes = {};

const mapStateToProps = ({ types }) => ({
  types
});

const mapDispatchToProps = (dispatch) => ({
  getAWSPricing: () => {
    dispatch(Actions.AWS.Pricing.getPricing())
  },
  getGCPPricing: () => {
    dispatch(Actions.GCP.Pricing.getPricing())
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(PageContainer);
