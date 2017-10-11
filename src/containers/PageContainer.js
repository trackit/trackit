import React, { Component } from 'react';
import {bindActionCreators} from 'redux';
import { connect } from 'react-redux';
// import PropTypes from 'prop-types';

import Paper from 'material-ui/Paper';


import DeclareSetupComponent from '../components/DeclareSetupComponent';
import ChartManagerComponent from '../components/ChartManagerComponent';
import TableComponent from '../components/TableComponent';


import * as ProvidersActions from '../actions/providersActions';

// PageContainer Component
class PageContainer extends Component {

    componentDidMount() {
      this.props.getPricingGcp();
      this.props.getPricingAws();
    }

    render() {
      const styles = {
        paper: {
          padding: '15px',
          margin: '15px',
        },
        noPadding: {
          padding: '0px'
        }
      }

      return (
        <div className="container-fluid">
          <Paper elevation={3} style={styles.paper} className="animated bounceInRight">
            <DeclareSetupComponent />
          </Paper>
          <Paper elevation={3} style={styles.paper} className="animated fadeIn">
            <ChartManagerComponent />
          </Paper>
          <Paper elevation={3} style={Object.assign({},styles.paper, styles.noPadding)} className="animated fadeIn">
            <TableComponent />
          </Paper>

        </div>
      );
    }

}


// Define PropTypes
PageContainer.propTypes = {};


// Subscribe component to redux store and merge the state into
// component's props
const mapStateToProps = ({ types }) => ({
  types
});

const mapActionCreatorsToProps = (dispatch) => (
   bindActionCreators(ProvidersActions, dispatch)
);


// connect method from react-router connects the component with redux store
export default connect(mapStateToProps, mapActionCreatorsToProps)(PageContainer);
