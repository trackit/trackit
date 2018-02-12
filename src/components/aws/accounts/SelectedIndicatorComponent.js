import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

// SelectedIndicator Component
export class SelectedIndicator extends Component {

  getText = () => {
    const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : '');

    if (!this.props.accounts.status)
      return null;
    if (this.props.accounts.status && (!this.props.accounts.values || !this.props.accounts.values.length || error))
      return `No AWS account available${error}`;
    if (this.props.selection.length === 0 || this.props.accounts.values.length === this.props.selection.length)
      return `${this.props.longVersion ? 'Displaying ' : ''}All accounts`;
    if (this.props.selection.length === 1)
      return `${this.props.longVersion ? 'Displaying ' : ''}${this.props.selection[0].pretty}`;
    return `${this.props.longVersion ? 'Displaying ' : ''}${this.props.selection.length} accounts`;
  };

  render() {
    const styles = {
      biggerBadge: {
        fontSize: '14px',
        fontWeight: '500',
      },
      icon: {
        fontSize: '16px',
      }
    };

    return(
      <span className="badge" style={styles.biggerBadge}>
        {this.props.icon && <span><i className="fa fa-amazon" style={styles.icon}/>&nbsp;&nbsp;</span>}
        {this.getText()}
      </span>
    );
  }

}

SelectedIndicator.propTypes = {
  accounts: PropTypes.shape({
    status: PropTypes.bool.isRequired,
    error: PropTypes.instanceOf(Error),
    values: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number.isRequired,
        roleArn: PropTypes.string.isRequired,
        pretty: PropTypes.string,
        bills: PropTypes.arrayOf(
          PropTypes.shape({
            bucket: PropTypes.string.isRequired,
            path: PropTypes.string.isRequired
          })
        ),
      })
    ),
  }),
  selection: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string,
    })
  ),
  longVersion: PropTypes.bool,
  icon: PropTypes.bool,
};

SelectedIndicator.defaultProps = {
  longVersion: false,
  icon: false,
};

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  selection: aws.accounts.selection,
});


export default connect(mapStateToProps)(SelectedIndicator);
