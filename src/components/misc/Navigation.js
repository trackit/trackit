import React, {Component} from 'react';
import { connect } from 'react-redux';
import { NavLink, withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import NavbarHeader from './NavbarHeader';
import Actions from '../../actions';

// Styling
import '../../styles/Navigation.css';

// Assets

export class Navigation extends Component {

  constructor() {
    super();
    this.state = {
      userMenuExpanded: false,
    };
  }

  toggleUserMenu() {
    this.setState({ userMenuExpanded: !this.state.userMenuExpanded });
  }

  render() {

    let userMenu;
    if (this.state.userMenuExpanded) {
      userMenu = (
        <ul className="nav nav-second-level">
          <li>
            <a href="" onClick={this.props.signOut}>
              <i className="menu-icon fa fa-sign-out"/>
              <span className="hide-menu">Sign out</span>
            </a>
          </li>
          <hr className="m-b-0"/>
        </ul>
      );
    }

    return(
      <div className="navigation">

        <NavbarHeader />

        <div className="navbar-default sidebar animated fadeInLeft" role="navigation">
          <div className="sidebar-head">
            <h3>
              <span className="open-close">
                <i className="fa fa-bars hidden-xs"></i>
                <i className="fa fa-times visible-xs"></i>
              </span>
              <span className="hide-menu">
                Navigation
              </span>
            </h3>
          </div>

          <ul className="nav" id="side-menu">
            <li className="user-menu-item">
              <button onClick={this.toggleUserMenu.bind(this)}>
                <span className="fa-stack fa-lg red-color">
                  <i className="fa fa-circle fa-stack-2x"></i>
                  <i className="fa fa-user fa-stack-1x fa-inverse"></i>
                </span>
                <span className="hide-menu">
                  {this.props.mail}
                  <i className="fa fa-caret-right"/>
                </span>
              </button>
              {userMenu || <hr className="m-b-0"/>}
            </li>

            <li>
              <NavLink exact to='/app' activeClassName="active">
                <i className="menu-icon fa fa-home"/>
                <span className="hide-menu">Home</span>
              </NavLink>
            </li>
            <li>
              <NavLink to='/app/s3' activeClassName="active">
                <i className="menu-icon fa fa-bar-chart"/>
                <span className="hide-menu">AWS S3 Analytics</span>
              </NavLink>
            </li>
            <li>
              <NavLink to='/app/optimize' activeClassName="active">
                <i className="menu-icon fa fa-area-chart"/>
                <span className="hide-menu">Compute Optimizer</span>
              </NavLink>
            </li>
            <li>
              <NavLink to='/app/ressources' activeClassName="active">
                <i className="menu-icon fa fa-pie-chart"/>
                <span className="hide-menu">Ressources Monitoring</span>
              </NavLink>
            </li>

            <li>
              <NavLink to='/app/setup' activeClassName="active">
                <i className="menu-icon fa fa-cog"/>
                <span className="hide-menu">Setup</span>
              </NavLink>
            </li>

          </ul>
        </div>
      </div>
    );
  }

}

Navigation.propTypes = {
  mail: PropTypes.string,
  signOut: PropTypes.func.isRequired
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  mail: state.auth.mail
});

/* istanbul ignore next */
const mapDispatchToProps = (dispatch) => ({
  signOut: () => {
    dispatch(Actions.Auth.logout())
  },
});

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Navigation));
