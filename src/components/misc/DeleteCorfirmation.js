import React, {Component} from 'react';
import Dialog from './Dialog';
import PropTypes from "prop-types";

class DeleteConfirmation extends Component {

  render() {
    return(
      <Dialog
        buttonName="Delete"
        buttonType="danger"
        title={"Are you sure to want to delete " + (this.props.entity || "this") + " ?"}
        actionName="Confirm"
        actionFunction={this.props.confirm}
      />
    );
  }

}

DeleteConfirmation.propTypes = {
  entity: PropTypes.string,
  confirm: PropTypes.func.isRequired
};


export default DeleteConfirmation;
