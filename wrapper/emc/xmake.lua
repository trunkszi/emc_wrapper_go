add_rules("mode.debug", "mode.release")
--add_rules("mode.debug")

add_requires("cocoyaxi", { configs = {
    vs_runtime = "MT",
}})

--add_requires("boost", { configs = {
--    vs_runtime = "MT",
--}})

-- xmake f --toolchain=mingw --sdk=C:\msys64\mingw64 -m debug
target("emc")
    set_arch("x64")
    set_kind("shared")
    set_languages("c++20")
    add_includedirs("src/third_party/gmsdk/include")
    if is_host("windows") then
        add_linkdirs("src/third_party/gmsdk/lib/win64")
        --add_ldflags("/DEFAULTLIB:libcmtd")
        --add_syslinks("libcmtd")
    end
    after_build(function(target)
    --    local pkgs = target:pkgs("cocoyaxi")["cocoyaxi"]["_INFO"]["libfiles"]
    --    for _, value in ipairs(pkgs)
    --    do
    --        os.cp(value, "$(buildir)/windows/x64/debug/")
    --    end
        os.cp("src/third_party/gmsdk/lib/win64/gmsdk.dll", "$(buildir)/windows/x64/debug/")
    end)
    add_files("src/interface.cpp")
    add_links("gmsdk")
    add_packages("cocoyaxi")

target("emc_demo")
    set_kind("binary")
    add_deps("emc")
    add_files("src/main.cpp")

--
-- If you want to known more usage about xmake, please see https://xmake.io
--
-- ## FAQ
--
-- You can enter the project directory firstly before building project.
--
--   $ cd projectdir
--
-- 1. How to build project?
--
--   $ xmake
--
-- 2. How to configure project?
--
--   $ xmake f -p [macosx|linux|iphoneos ..] -a [x86_64|i386|arm64 ..] -m [debug|release]
--
-- 3. Where is the build output directory?
--
--   The default output directory is `./build` and you can configure the output directory.
--
--   $ xmake f -o outputdir
--   $ xmake
--
-- 4. How to run and debug target after building project?
--
--   $ xmake run [targetname]
--   $ xmake run -d [targetname]
--
-- 5. How to install target to the system directory or other output directory?
--
--   $ xmake install
--   $ xmake install -o installdir
--
-- 6. Add some frequently-used compilation flags in xmake.lua
--
-- @code
--    -- add debug and release modes
--    add_rules("mode.debug", "mode.release")
--
--    -- add macro defination
--    add_defines("NDEBUG", "_GNU_SOURCE=1")
--
--    -- set warning all as error
--    set_warnings("all", "error")
--
--    -- set language: c99, c++11
--    set_languages("c99", "c++11")
--
--    -- set optimization: none, faster, fastest, smallest
--    set_optimize("fastest")
--
--    -- add include search directories
--    add_includedirs("/usr/include", "/usr/local/include")
--
--    -- add link libraries and search directories
--    add_links("tbox")
--    add_linkdirs("/usr/local/lib", "/usr/lib")
--
--    -- add system link libraries
--    add_syslinks("z", "pthread")
--
--    -- add compilation and link flags
--    add_cxflags("-stdnolib", "-fno-strict-aliasing")
--    add_ldflags("-L/usr/local/lib", "-lpthread", {force = true})
--
-- @endcode
--

