import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

// SelectedIndicator Component
class SelectedIndicator extends Component {

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

    const getText = () => {
      const error = (this.props.accounts.error ? ` (${this.props.accounts.error.message})` : null);
      if (this.props.accounts.status && (!this.props.accounts.values || !this.props.accounts.values.length || error))
        return `No AWS account available${error}`;
      if (this.props.selection.length === 0)
        return `${this.props.longVersion ? 'Displaying' : ''} All accounts`;
      if (this.props.selection.length === 1)
        return `${this.props.longVersion ? 'Displaying' : ''} ${this.props.selection[0].pretty}`;
      return `${this.props.longVersion ? 'Displaying' : ''} ${this.props.selection.length} accounts`;
    };

    return(
      <span className="badge" style={styles.biggerBadge}>
        {this.props.icon && <span><i className="fa fa-amazon" style={styles.icon}/>&nbsp;&nbsp;</span>}
        {getText()}
      </span>
    );
  }

}

SelectedIndicator.defaultProps = {
  longVersion: false,
  icon: false,
};

SelectedIndicator.propTypes = {
  accounts: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      roleArn: PropTypes.string.isRequired,
      pretty: PropTypes.string,
    })
  ),
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

/* istanbul ignore next */
const mapStateToProps = ({aws}) => ({
  accounts: aws.accounts.all,
  selection: aws.accounts.selection,
});


export default connect(mapStateToProps)(SelectedIndicator);
