import React, {Component} from 'react';
import { connect } from 'react-redux';
import { NavLink, withRouter } from 'react-router-dom';
import PropTypes from 'prop-types';
import NavbarHeader from './NavbarHeader';
import AccountSelector from '../aws/accounts/SelectorComponent';
import SelectedIndicator from '../aws/accounts/SelectedIndicatorComponent';

// Styling
import '../../styles/Navigation.css';

// Assets

export class Navigation extends Component {

  constructor() {
    super();
    this.state = {
      userMenuExpanded: false,
    };
    this.closeUserMenu = this.closeUserMenu.bind(this);
  }

  toggleUserMenu() {
    this.setState({ userMenuExpanded: !this.state.userMenuExpanded });
  }

  closeUserMenu = (e) => {
    e.preventDefault();
    this.setState({ userMenuExpanded: false });
  };

  render() {

    let userMenu;
    if (this.state.userMenuExpanded) {
      userMenu = (
        <div className="nav nav-second-level">
          <AccountSelector/>
        </div>
      );
    }

    return(
      <div className="navigation">

        <NavbarHeader />

        <div className="navbar-default sidebar animated fadeInLeft" role="navigation" onMouseLeave={this.closeUserMenu}>
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
                  <i className="fa fa-amazon fa-stack-1x fa-inverse"></i>
                </span>
                <span className="hide-menu">
                  <SelectedIndicator />
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
            {/*
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
            */}
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
  mail: PropTypes.string
};

/* istanbul ignore next */
const mapStateToProps = (state) => ({
  mail: state.auth.mail
});

export default withRouter(connect(mapStateToProps)(Navigation));
