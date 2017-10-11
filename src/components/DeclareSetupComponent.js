import React, { Component } from 'react';
import {bindActionCreators} from 'redux';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import Button from 'material-ui/Button';

import * as ProvidersActions from '../actions/providersActions';

// DeclareSetupComponent Component
class DeclareSetupComponent extends Component {

    submit() {
      const payload = {
        frequentValue: this.frequentInput.value,
        frequentUnit: this.frequentSelect.value,
        infrequentValue: this.infrequentInput.value,
        infrequentUnit: this.infrequentSelect.value,
        archiveValue: this.archiveInput.value,
        archiveUnit: this.archiveSelect.value,
      };
      this.props.setStorageTypes(payload);
      this.props.getPricingGcp();
      this.props.getPricingAws();
    }

    render() {
      const { types } = this.props;

      const styles = {
        titles: {
          textAlign: 'left',
        },
        inputs: {
          textAlign: 'right',
          width: '50%',
        },
        selects: {
          width: '60px',
        },
        submitBtn: {
          width: '30%',
          color: 'white',
          backgroundColor: '#d6413b',
          margin: '15px auto',
          display: 'block',
        },
      };

      return (
        <div>
          <div>
            <div className="col-md-12">
              <h4 className="paper-title">
                <i className="fa fa-cog red-color"/>
                &nbsp;
                Setup
              </h4>
            </div>
          </div>
          <div className="clearfix" />
          <hr />
          <div className="row">
            <div className="col-md-4 form-inline">
              <h4 style={styles.titles}>Frequent Access</h4>
              <input
                className="form-control pull-left"
                type="number"
                defaultValue={types.frequentValue}
                style={styles.inputs}
                ref={(input) => this.frequentInput = input}
              />
              &nbsp;
              <select
                className="form-control pull-left"
                style={styles.selects}
                defaultValue={types.frequentUnit}
                ref={(input) => this.frequentSelect = input}
              >
                <option value={'GB'}>GB</option>
                <option value={'TB'}>TB</option>
              </select>
              <div style={{ clear: 'both' }} />
            </div>

            <div className="col-md-4 form-inline">
              <h4 style={styles.titles}>Infrequent Access</h4>
              <input
                className="form-control pull-left"
                type="number"
                style={styles.inputs}
                defaultValue={types.infrequentValue}
                ref={(input) => this.infrequentInput = input}
              />
                &nbsp;
              <select
                className="form-control pull-left"
                style={styles.selects}
                defaultValue={types.infrequentUnit}
                ref={(input) => this.infrequentSelect = input}
              >
                <option value={'GB'}>GB</option>
                <option value={'TB'}>TB</option>
              </select>
              <div style={{ clear: 'both' }} />
            </div>


            <div className="col-md-4 form-inline">
              <h4 style={styles.titles}>Archive</h4>
              <input
                className="form-control pull-left"
                type="number"
                style={styles.inputs}
                defaultValue={types.archiveValue}
                ref={(input) => this.archiveInput = input}
              />
              &nbsp;
              <select
                className="form-control pull-left"
                style={styles.selects}
                defaultValue={types.archiveUnit}
                ref={(input) => this.archiveSelect = input}
              >
                <option value={'GB'}>GB</option>
                <option value={'TB'}>TB</option>
              </select>
              <div style={{ clear: 'both' }} />
            </div>

            <div style={{ clear: 'both' }} />

            <Button
              style={styles.submitBtn}
              onClick={this.submit.bind(this)}
              color="primary"
              raised={true}
            >
              <i className="fa fa-refresh" />
              &nbsp;
              Refresh
            </Button>

          </div>
        </div>
      );
    }

}

// Define PropTypes
DeclareSetupComponent.propTypes = {
  types: PropTypes.object,
};

// Subscribe component to redux store and merge the state into
// component's props
const mapStateToProps = ({ types }) => ({
  types
});

const mapActionCreatorsToProps = (dispatch) => (
   bindActionCreators(ProvidersActions, dispatch)
);


// connect method from react-router connects the component with redux store
export default connect(mapStateToProps, mapActionCreatorsToProps)(DeclareSetupComponent);
